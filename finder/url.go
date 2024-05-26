package finder

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
	urlFmt         = "%s%s?classificationName=%s&city=%s&stateCode=%s&radius=%s&unit=%s&localStartDateTime=%s&sort=%s&size=%v"
	apiKeyFmt      = "&apikey=%s"
	dateTimeFmt    = "2006-01-02T15:04:05"
	dateFmt        = "2006-01-02"
	classification = "music"
	sort           = "date,asc"
	radius         = "50"
	unit           = "miles"
	stateCode      = "GA"
	pageSize       = 20
)

type Url string
type UrlPath string

func buildTicketmasterUrl(city string, state string) (Url, error) {
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

	url := fmt.Sprintf(urlFmt, host, eventPath, classification, city, state, radius, unit, startDate, sort, pageSize)
	log.Debug("Built URL (without auth token): ", url)
	url += fmt.Sprintf(apiKeyFmt, token)
	return Url(url), nil
}

func buildTicketmasterUrlWithPath(path UrlPath) (Url, error) {
	token, err := getAuthToken()
	if err != nil {
		return "", err
	}

    url := string(host + path)
	log.Debug("Built URL (without auth token): ", url)
	url += fmt.Sprintf(apiKeyFmt, token)
	return Url(url), nil
}

func getAuthToken() (string, error) {
	token := os.Getenv(apiKey)
	if token == "" {
		errMsg := fmt.Sprintf("%s environment variable must be set", apiKey)
		return "", errors.New(errMsg)
	}
	return token, nil
}
