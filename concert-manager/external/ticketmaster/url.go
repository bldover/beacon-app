package ticketmaster

import (
	"concert-manager/log"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	apiKey      = "CM_TICKETMASTER_API_KEY"
	host        = "https://app.ticketmaster.com"
	eventPath   = "/discovery/v2/events"
	urlFmt      = "%s%s?classificationName=music&geoPoint=%s&radius=%d&unit=miles&localStartDateTime=%s&sort=%s&size=%v"
	apiKeyFmt   = "&apikey=%s"
	dateTimeFmt = "2006-01-02T15:04:05"
	dateFmt     = "2006-01-02"
	sort        = "date,asc"
	radius      = 50
	pageSize    = 50
)

func buildTicketmasterUrl(city, stateCd string) (string, error) {
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

	geopoint := getGeopoint(city, stateCd)

	url := fmt.Sprintf(urlFmt, host, eventPath, geopoint, radius, startDate, sort, pageSize)
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

func getGeopoint(city, stateCd string) string {
	geopointMap := map[string]string{
		"atlanta-ga": "dn5bzz", // centered around sandy springs
	}
	key := strings.ToLower(city) + "-" + strings.ToLower(stateCd)
	return geopointMap[key]
}
