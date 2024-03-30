package finder

import (
	"concert-manager/data"
	"concert-manager/log"
	"concert-manager/util"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	quotaViolationCode = "policies.ratelimit.QuotaViolation"
	rateViolationCode  = "policies.ratelimit.SpikeArrestViolation"
)

type RetryableError struct {
	message string
}

func (e RetryableError) Error() string {
	return e.message
}

type EventCancelledError struct {
	message string
}

func (e EventCancelledError) Error() string {
	return e.message
}

type EventCount struct {
    successCount int
	cancelledCount int
	failedCount int
}

type ticketmasterRetriever struct {}

func (r ticketmasterRetriever) GetUpcomingEvents(request FindEventRequest) ([]data.EventDetails, error) {
	city := request.City
	state := request.State
	log.Infof("Starting to retrieve all upcoming events from Ticketmaster for %s, %s", city, state)

	url, err := buildTicketmasterUrl(city, state)
	if err != nil {
		return nil, err
	}
	response, err := getResponseDetails(url)
	if err != nil {
		// Assume no rate violation here since it's the first request
		log.Error("Error retrieving event data from Ticketmaster", err)
		return nil, err
	}

	expectedEventCount := response.PageInfo.EventCount
	eventDetails := make([]data.EventDetails, 0, expectedEventCount)

	eventCount, err := populateAllEventDetails(response, &eventDetails)
	if err != nil {
		log.Error(err)
	}

	nextUrlPath := UrlPath(response.Links.Next.URL)
	remainingEventCount := getRemainingPages(nextUrlPath, &eventDetails)

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

func getRemainingPages(urlPath UrlPath, eventDetails *[]data.EventDetails) EventCount {
	retryCount := 0
	maxRetries := 3
	eventCount := EventCount{}
	for urlPath != "" {
		// try to not exceed the rate limit
		waitTime := time.Duration(200000000) // 0.2s
		time.Sleep(waitTime)
		lastUrlPath := urlPath
		var err error
		var pageEventCount EventCount
		urlPath, pageEventCount, err = getEvents(urlPath, eventDetails)
		if err != nil {
			switch err.(type) {
			case RetryableError:
				if retryCount < maxRetries {
					log.Info("Received Ticketmaster rate violation, retry count:", retryCount)
					urlPath = lastUrlPath
					retryCount++
					waitTime := time.Duration(500000000) // 0.5s
					time.Sleep(waitTime)
					continue
				} else {
					log.Error("Failed to retrieve event page from Ticketmaster after all retry attempts:", err)
					eventCount.failedCount += pageEventCount.failedCount
					break
				}
			case error:
				log.Error("Failed to retrieve event page from Ticketmaster with non-retryable error:", err)
				eventCount.failedCount += pageEventCount.failedCount
				break
			}
		}
		log.Debug("Successfully retrieved event page from Ticketmaster")
		eventCount.successCount += pageEventCount.successCount
		eventCount.cancelledCount += pageEventCount.cancelledCount
		eventCount.failedCount += pageEventCount.failedCount
		retryCount = 0
	}
	return eventCount
}

func getEvents(urlPath UrlPath, events *[]data.EventDetails) (UrlPath, EventCount, error) {
	eventCount := EventCount{}
	url, err := buildTicketmasterUrlWithPath(urlPath)
	if err != nil {
		eventCount.failedCount += pageSize
		return "", eventCount, err
	}

	response, err := getResponseDetails(url)
	if err != nil {
		eventCount.failedCount += pageSize
		return "", eventCount, err
	}
	pageEventCount, err := populateAllEventDetails(response, events);
	if err != nil {
		log.Error(err)
	}

	eventCount.successCount += pageEventCount.successCount
	eventCount.cancelledCount += pageEventCount.cancelledCount
	eventCount.failedCount += pageEventCount.failedCount
	return UrlPath(response.Links.Next.URL), eventCount, nil
}

func getResponseDetails(url Url) (*tmResponse, error) {
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
			return nil, RetryableError{"exceeded per-second rate limit"}
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

func populateAllEventDetails(response *tmResponse, events *[]data.EventDetails) (EventCount, error) {
	eventCount := EventCount{}
	for _, event := range response.Data.Events {
		eventDetails, err := parseEventDetails(&event)
		if err != nil {
			switch err.(type) {
			case EventCancelledError:
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

func parseEventDetails(event *tmEventResponse) (*data.EventDetails, error) {
	eventName := event.EventName
	artistDetails := event.Details.Artists
	if eventName == "" && len(artistDetails) == 0 {
		return nil, errors.New("no event name or artists")
	}

	mainAct := data.Artist{}
	if len(artistDetails) != 0 {
		mainActDetails := artistDetails[0]
		mainAct.Name = mainActDetails.Name
		if len(mainActDetails.Classification) != 0 {
			mainAct.Genre = getGenre(mainActDetails.Classification[0])
		}
	}

	openers := []data.Artist{}
	if len(artistDetails) > 1 {
		for _, openerDetails := range artistDetails[1:] {
			if openerDetails.Name == "" {
				return nil, errors.New("no opener artist name")
			}
			opener := data.Artist{
				Name: openerDetails.Name,
			}
			if len(openerDetails.Classification) != 0 {
				opener.Genre = getGenre(openerDetails.Classification[0])
			}
			openers = append(openers, opener)
		}
	}

	eventGenre := ""
	if len(event.Classification) != 0 {
		eventGenre = getGenre(event.Classification[0])
	}

	price := ""
	if len(event.Prices) == 0 {
		price = "Unknown"
	} else {
 		price = strconv.FormatFloat(event.Prices[0].MinPrice, 'f', 2, 64)
		if !event.Ticketing.InclusivePricing.Enabled {
			price += " + fees"
		}
	}

	venue := data.Venue{}
	if len(event.Details.Venues) != 0 {
		venueDetails := event.Details.Venues[0]
		venue.Name = venueDetails.Name
		venue.City = venueDetails.City.Name
		venue.State = venueDetails.State.Name
	}

	dateRaw := event.Dates.Start.Date
	date, err := time.Parse(dateFmt, dateRaw)
	if err != nil {
		errMsg := fmt.Sprintf("unable to parse event date %s", dateRaw)
		return nil, errors.New(errMsg)
	}

	eventDetails := data.EventDetails{
		Name:  eventName,
		Price: price,
		EventGenre: eventGenre,
		Event: data.Event{
			MainAct: mainAct,
			Openers: openers,
			Venue:   venue,
			Date:    util.Date(date),
		},
	}

	if event.Dates.Status.Code == "cancelled" {
		return &eventDetails, EventCancelledError{"Event has been cancelled"}
	}
	return &eventDetails, nil
}

func getGenre(genres tmGenreResponse) string {
	subGenre := genres.Subgenre.Name
	genre := genres.Genre.Name
	switch {
	case subGenre != "" && subGenre != "Undefined" && subGenre != "Other":
		return subGenre
	case genre != "" && genre != "Undefined" && genre != "Other":
		return genre
	default:
		return ""
	}
}
