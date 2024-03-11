package ticketmaster

import (
	"concert-manager/data"
	"concert-manager/log"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	apiKey         = "TICKETMASTER_API_KEY"
	eventsBaseUrl  = "https://app.ticketmaster.com/discovery/v2/events"
	urlFmt         = "%s?classificationName=%s&city=%s&stateCode=%s&radius=%s&unit=%s&localStartDateTime=%s&sort=%s&size=%s"
	apiKeyFmt      = "&apikey=%s"
	dateTimeFmt    = "2006-01-02T15:04:05"
	dateFmt        = "2006-01-02"
	classification = "music"
	sort           = "date,asc"
	radius         = "50"
	unit           = "miles"
	stateCode      = "GA"
	size           = "20"

	quotaViolationCode = "policies.ratelimit.QuotaViolation"
	rateViolationCode  = "policies.ratelimit.SpikeArrestViolation"
)

type response struct {
	Links struct {
		Next struct {
			URL string `json:"href"`
		} `json:"next"`
	} `json:"_links"`
	Data struct {
		Events []eventResponse `json:"events"`
	} `json:"_embedded"`
	PageInfo struct {
		EventCount int `json:"totalElements"`
	} `json:"page"`
}

type eventResponse struct {
	EventName string `json:"name"`
	Dates     struct {
		Start struct {
			Date string `json:"localDate"`
		} `json:"start"`
	} `json:"dates"`
	Prices []struct {
		MinPrice float64 `json:"min"`
	} `json:"priceRanges"`
	Ticketing struct {
		InclusivePricing struct {
			Enabled bool `json:"enabled"`
		} `json:"allInclusivePricing"`
	} `json:"ticketing"`
	Details struct {
		Venues []struct {
			Name string `json:"name"`
			City struct {
				Name string `json:"Name"`
			} `json:"city"`
			State struct {
				Name string `json:"name"`
			} `json:"state"`
		} `json:"venues"`
		Artists []struct {
			Name  string `json:"name"`
			Links struct {
				Wiki []struct {
					URL string `json:"url"`
				} `json:"wiki"`
				Spotify []struct {
					URL string `json:"url"`
				} `json:"spotify"`
			} `json:"externalLinks"`
			Classification []struct {
				Genre struct {
					Name string `json:"name"`
				} `json:"genre"`
				Subgenre struct {
					Name string `json:"name"`
				} `json:"subGenre"`
			} `json:"classifications"`
		} `json:"attractions"`
	} `json:"_embedded"`
}

type errorResponse struct {
	Fault struct {
		Details struct {
			Code string `json:"errorcode"`
		} `json:"detail"`
	} `json:"fault"`
}

type RetryableError struct {
	message string
}

func (e RetryableError) Error() string {
	return e.message
}

type UpcomingEventsRequest struct {
    city string
	state string
}

//func (r *UpcomingEventsRequest)

// This function is the entry point to getting all the event details.
// It retrieves the data from the first page of event results and then submits a
// bunch of parallel calls to getEventsByPage to retrieve all the other page data
func GetUpcomingEvents(city string, state string) (*[]*data.EventDetails, error) {
	log.Infof("Starting to retrieve all upcoming events from Ticketmaster for %s, %s", city, state)
	url, err := buildUrl(city, state)
	if err != nil {
		return nil, err
	}
	response, err := getResponseDetails(url)
	if err != nil {
		// Don't check for retryable error because this is the first request,
		// so we assume we haven't exceeded the rate limit.
		// This will be inappropriate and need to be changed if the app ever
		// reaches a point where this method would be called this simultaneously
		// by multiple users or some other logic
		return nil, err
	}

	eventCount := response.PageInfo.EventCount
	eventDetails := make([]*data.EventDetails, eventCount)
	populateAllEventDetails(response, &eventDetails)
	// handle error scenario

//	doneIndicator := make(chan int)
	// submit next call if page exists
	return &eventDetails, nil
}

func getEventsByPage(url string, events *[]*data.EventDetails, done chan int) error {
	token, err := getAuthToken()
	if err != nil {
		return err
	}
	url += fmt.Sprintf("%s&apiKey=%s", url, token)

	response, err := getResponseDetails(url)
	populateAllEventDetails(response, events)
	// handle error scenario

	// submit next call if page exists
	done <- 1
	return nil
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

func buildUrl(city string, state string) (string, error) {
	token, err := getAuthToken()
	if err != nil {
		return "", err
	}

	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		errMsg := fmt.Sprintf("failed to find time zone with err: %v", err)
		return "", errors.New(errMsg)
	}
	startDate := time.Now().In(location).Format(dateTimeFmt)

	url := fmt.Sprintf(urlFmt, eventsBaseUrl, classification, city, state, radius, unit, startDate, sort, size)
	log.Debug("Built URL (without auth token): ", url)
	url += fmt.Sprintf(apiKeyFmt, token)
	log.Debug("Built full URL: ", url)
	return url, nil
}

func getAuthToken() (string, error) {
    token := os.Getenv(apiKey)
	if token == "" {
		errMsg := fmt.Sprintf("%s environment variable must be set", apiKey)
		return "", errors.New(errMsg)
	}
	return token, nil
}

func toResponse(body io.Reader) (*response, error) {
	var resp response
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		errMsg := fmt.Sprintf("failed to parse ticketmaster response: %v", err)
		return nil, errors.New(errMsg)
	}
	return &resp, nil
}

func toErrorResponse(body io.Reader) (*errorResponse, error) {
	var resp errorResponse
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		errMsg := fmt.Sprintf("failed to parse ticketmaster error response: %v", err)
		return nil, errors.New(errMsg)
	}
	return &resp, nil
}

func populateAllEventDetails(response *response, events *[]*data.EventDetails) error {
	eventCount := 0
	for _, event := range response.Data.Events {
		eventDetails, err := parseEventDetails(&event)
		if err != nil {
			log.Errorf("Failed to parse event %+v, with error %v", event, err)
			continue
		}
		*events = append(*events, eventDetails)
		eventCount++
	}

	if eventCount != len(response.Data.Events) {
		return errors.New("failed to parse all events")
	}
	return nil
}

func parseEventDetails(event *eventResponse) (*data.EventDetails, error) {
	name := event.EventName
	if name == "" {
		return nil, errors.New("no event name")
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

	if len(event.Details.Venues) == 0 {
		return nil, errors.New("no venue")
	}
	venueDetails := event.Details.Venues[0]
	venue := data.Venue{
		Name:  venueDetails.Name,
		City:  venueDetails.City.Name,
		State: venueDetails.State.Name,
	}
	if !venue.Populated() {
		return nil, errors.New("missing some venue data")
	}

	if len(event.Details.Artists) == 0 {
		return nil, errors.New("no artists")
	}
	artistDetails := event.Details.Artists
	mainActDetails := artistDetails[0]
	if mainActDetails.Name == "" {
		return nil, errors.New("no main act artist name")
	}
	mainAct := data.Artist{
		Name: mainActDetails.Name,
	}
	if len(mainActDetails.Classification) != 0 {
		mainAct.Genre = mainActDetails.Classification[0].Subgenre.Name
	}

	openers := []data.Artist{}
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

	dateRaw := event.Dates.Start.Date
	date, err := time.Parse(dateFmt, dateRaw)
	if err != nil {
		errMsg := fmt.Sprintf("unable to parse event date %s", dateRaw)
		return nil, errors.New(errMsg)
	}

	eventDetails := data.EventDetails{
		Name:  name,
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
