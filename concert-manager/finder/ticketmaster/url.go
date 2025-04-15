package ticketmaster

import (
	"concert-manager/log"
	"errors"
	"fmt"
	"os"
	"time"
)

const (
	apiKey         = "CM_TICKETMASTER_API_KEY"
	host           = "https://app.ticketmaster.com"
	eventPath      = "/discovery/v2/events"
	// radius/unit doesn't seem to work with city/state, look into getting latlong?
	urlFmt         = "%s%s?classificationName=%s&stateCode=%s&radius=%s&unit=%s&localStartDateTime=%s&sort=%s&size=%v"
	apiKeyFmt      = "&apikey=%s"
	dateTimeFmt    = "2006-01-02T15:04:05"
	dateFmt        = "2006-01-02"
	classification = "music"
	sort           = "date,asc"
	radius         = "50"
	unit           = "miles"
	stateCode      = "GA"
	pageSize       = 50
)

func buildTicketmasterUrl(state string) (string, error) {
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

	url := fmt.Sprintf(urlFmt, host, eventPath, classification, state, radius, unit, startDate, sort, pageSize)
	log.Debug("Built URL (without auth token): ", url)
	url += fmt.Sprintf(apiKeyFmt, token)
	return url, nil
}

func buildTicketmasterUrlWithPath(path string) (string, error) {
	token, err := getAuthToken()
	if err != nil {
		return "", err
	}

    url := string(host + path)
	log.Debug("Built URL (without auth token): ", url)
	url += fmt.Sprintf(apiKeyFmt, token)
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
