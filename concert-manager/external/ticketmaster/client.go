package ticketmaster

import (
	"concert-manager/domain"
	"concert-manager/log"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var TEST_MODE = false

const testModePageLimit = 3

const (
	quotaViolationCode = "policies.ratelimit.QuotaViolation"
	rateViolationCode  = "policies.ratelimit.SpikeArrestViolation"
)

type (
	retryableError struct {
		message string
	}
	eventCancelledError struct {
		message string
	}
)

func (e retryableError) Error() string {
	return e.message
}

func (e eventCancelledError) Error() string {
	return e.message
}

type eventCount struct {
	successCount   int
	cancelledCount int
	failedCount    int
}

type Ticketmaster struct{}

func (t Ticketmaster) GetUpcomingEvents(city string, stateCd string) ([]domain.EventDetails, error) {
	log.Infof("Starting to retrieve all upcoming events from Ticketmaster for %s", stateCd)

	url, err := buildTicketmasterUrl(city, stateCd)
	if err != nil {
		return nil, err
	}
	response, err := t.getResponseDetails(url)
	if err != nil {
		// Assume no rate violation here since it's the first request
		log.Error("Error retrieving event data from Ticketmaster", err)
		return nil, err
	}

	// ticketmaster max 1k events
	expectedEventCount := min(response.PageInfo.EventCount, 1000)
	eventDetails := make([]domain.EventDetails, 0, expectedEventCount)

	eventCount, err := t.populateAllEventDetails(response, &eventDetails)
	if err != nil {
		log.Error(err)
	}

	nextUrlPath := response.Links.Next.URL
	remainingEventCount := t.getRemainingPages(nextUrlPath, &eventDetails)

	eventCount.successCount += remainingEventCount.successCount
	eventCount.failedCount += remainingEventCount.failedCount
	eventCount.cancelledCount += remainingEventCount.cancelledCount

	expectedNotCancelledCount := expectedEventCount - eventCount.cancelledCount
	if len(eventDetails) != expectedNotCancelledCount {
		errFmt := "Unable to retrieve all expected events. Read %v/%v"
		errMsg := fmt.Sprintf(errFmt, len(eventDetails), expectedNotCancelledCount)
		return eventDetails, errors.New(errMsg)
	}

	log.Infof("Ticketmaster read counts: %+v", eventCount)
	return eventDetails, nil
}

func (t Ticketmaster) getRemainingPages(urlPath string, eventDetails *[]domain.EventDetails) eventCount {
	retryCount := 0
	maxRetries := 3
	count := eventCount{}
	total := 0
	// ticketmaster restricts paging to the first 1k events only
	// subtract since the first page of results isn't included
	for urlPath != "" && total < (1000-pageSize) {
		// try to not exceed the 5/s rate limit
		waitTime := 200 * time.Millisecond
		time.Sleep(waitTime)
		lastUrlPath := urlPath
		var err error
		var pageEventCount eventCount
		urlPath, pageEventCount, err = t.getEvents(urlPath, eventDetails)
		if err != nil {
			switch err.(type) {
			case retryableError:
				if retryCount < maxRetries {
					log.Info("Received Ticketmaster rate violation, retry count:", retryCount)
					urlPath = lastUrlPath
					retryCount++
					waitTime := 500 * time.Millisecond
					time.Sleep(waitTime)
					continue
				} else {
					log.Error("Failed to retrieve event page from Ticketmaster after all retry attempts:", err)
					count.failedCount += pageEventCount.failedCount
					break
				}
			case error:
				log.Error("Failed to retrieve event page from Ticketmaster with non-retryable error:", err)
				count.failedCount += pageEventCount.failedCount
			}
		}
		log.Debug("Successfully retrieved event page from Ticketmaster")
		count.successCount += pageEventCount.successCount
		count.cancelledCount += pageEventCount.cancelledCount
		count.failedCount += pageEventCount.failedCount
		total += pageSize
		retryCount = 0
	}
	return count
}

func (t Ticketmaster) getEvents(urlPath string, events *[]domain.EventDetails) (string, eventCount, error) {
	eventCount := eventCount{}
	url, err := buildTicketmasterUrlWithPath(urlPath)
	if err != nil {
		eventCount.failedCount += pageSize
		return "", eventCount, err
	}

	response, err := t.getResponseDetails(url)
	if err != nil {
		eventCount.failedCount += pageSize
		return "", eventCount, err
	}
	log.Debugf("Received ticketmaster response page: %+v", *response)
	pageEventCount, err := t.populateAllEventDetails(response, events)
	if err != nil {
		log.Error(err)
	}

	eventCount.successCount += pageEventCount.successCount
	eventCount.cancelledCount += pageEventCount.cancelledCount
	eventCount.failedCount += pageEventCount.failedCount

	if TEST_MODE && response.PageInfo.Page >= testModePageLimit {
		return "", eventCount, nil
	}
	return response.Links.Next.URL, eventCount, nil
}

func (t Ticketmaster) getResponseDetails(url string) (*tmResponse, error) {
	response, err := http.Get(string(url))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errResp, err := toErrorResponse(response.Body)
		if err != nil {
			return nil, err
		}

		if response.StatusCode == http.StatusTooManyRequests && errResp.Fault.Details.Code == rateViolationCode {
			return nil, retryableError{"exceeded per-second rate limit"}
		}

		errFmt := "received error code %v: %s from ticketmaster with details %v"
		errMsg := fmt.Sprintf(errFmt, response.StatusCode, response.Status, errResp)
		return nil, errors.New(errMsg)
	}

	respData, err := toResponse(response.Body)
	if err != nil {
		return nil, err
	}
	return respData, nil
}

func (t Ticketmaster) populateAllEventDetails(response *tmResponse, events *[]domain.EventDetails) (eventCount, error) {
	eventCount := eventCount{}
	for _, event := range response.Data.Events {
		eventDetails, err := parseEventDetails(&event)
		if err != nil {
			switch err.(type) {
			case eventCancelledError:
				log.Debugf("Skipped event %+v due to being cancelled", eventDetails)
				eventCount.cancelledCount++
				continue
			case error:
				log.Errorf("Failed to parse event %+v, with error %v", event, err)
				eventCount.failedCount++
				continue
			}
		}
		*events = append(*events, *eventDetails)
		eventCount.successCount++
	}

	if eventCount.failedCount != 0 {
		return eventCount, fmt.Errorf("failed to parse %v events", eventCount.failedCount)
	}
	return eventCount, nil
}
