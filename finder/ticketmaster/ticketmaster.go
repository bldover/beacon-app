package ticketmaster

import (
	"concert-manager/data"
	"concert-manager/log"
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

type UpcomingEventsRequest struct {
	City  string
	State string
}

// This function is the entry point to getting all the event details.
// It retrieves the data from the first page of event results and then
// starts a chain of calls to retrieve the data from the other pages
func GetUpcomingEvents(request UpcomingEventsRequest) ([]data.EventDetails, error) {
	city := request.City
	state := request.State
	log.Infof("Starting to retrieve all upcoming events from Ticketmaster for %s, %s", city, state)

	url, err := buildUrl(city, state)
	if err != nil {
		return nil, err
	}
	response, err := getResponseDetails(url)
	if err != nil {
		// Assume no rate violation here since it's the first request
		log.Error("Error retrieving event data from Ticketmaster", err)
		return nil, err
	}

	eventCount := response.PageInfo.EventCount
	eventDetails := make([]data.EventDetails, 0, eventCount)
	if err := populateAllEventDetails(response, &eventDetails); err != nil {
		log.Error(err)
	}

	getRemainingPages(response.Links.Next.URL, &eventDetails)
	if len(eventDetails) != eventCount {
		errMsg := fmt.Sprintf("Unable to retrieve all expected events. Read %v/%v", len(eventDetails), eventCount)
		return eventDetails, errors.New(errMsg)
	}
	return eventDetails, nil
}

func getRemainingPages(url string, eventDetails *[]data.EventDetails) {
	retryCount := 0
	maxRetries := 3
	for url != "" {
		// try to not exceed the rate limit
		waitTime := time.Duration(200000000) // 0.2s
		time.Sleep(waitTime)
		lastUrl := url
		var err error
		url, err = getEvents(url, eventDetails)
		if err != nil {
			switch err.(type) {
			case RetryableError:
				if retryCount < maxRetries {
					log.Info("Received Ticketmaster rate violation, retry count:", retryCount)
					url = lastUrl
					retryCount++
					waitTime := time.Duration(500000000) // 0.5s
					time.Sleep(waitTime)
				} else {
					log.Error("Failed to retrieve event page from Ticketmaster after all retry attempts:", err)
				}
			case error:
				log.Error("Failed to retrieve event page from Ticketmaster with non-retryable error:", err)
			}
			continue
		}
		log.Debug("Successfully retrieved event page from Ticketmaster")
		retryCount = 0
	}
}

func getEvents(urlPath string, events *[]data.EventDetails) (string, error) {
	token, err := getAuthToken()
	if err != nil {
		return "", err
	}
	url := host + urlPath
	log.Debug("Built URL (without auth token): ", url)
	url += fmt.Sprintf(apiKeyFmt, token)

	response, err := getResponseDetails(url)
	if err != nil {
		return "", err
	}
	if err := populateAllEventDetails(response, events); err != nil {
		log.Error(err)
	}

	return response.Links.Next.URL, nil
}

func getResponseDetails(url string) (*response, error) {
	response, err := http.Get(url)
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

func populateAllEventDetails(response *response, events *[]data.EventDetails) error {
	eventCount := 0
	for _, event := range response.Data.Events {
		eventDetails, err := parseEventDetails(&event)
		if err != nil {
			log.Errorf("Failed to parse event %+v, with error %v", event, err)
			continue
		}
		*events = append(*events, *eventDetails)
		eventCount++
	}

	failedCount := len(response.Data.Events) - eventCount
	if failedCount != 0 {
		return fmt.Errorf("failed to parse %v events", failedCount)
	}
	return nil
}

func parseEventDetails(event *eventResponse) (*data.EventDetails, error) {
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
			mainAct.Genre = mainActDetails.Classification[0].Subgenre.Name
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
				opener.Genre = openerDetails.Classification[0].Subgenre.Name
			}
			openers = append(openers, opener)
		}
	}

	price := ""
	if len(event.Prices) == 0 {
		price = "unknown"
	} else {
 		price = strconv.FormatFloat(event.Prices[0].MinPrice, 'f', 2, 64)
	}
	if !event.Ticketing.InclusivePricing.Enabled {
		price += " + fees"
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
		Event: data.Event{
			MainAct: mainAct,
			Openers: openers,
			Venue:   venue,
			Date:    data.Date(date),
		},
	}
	return &eventDetails, nil
}
