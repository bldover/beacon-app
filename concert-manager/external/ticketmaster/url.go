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
	apiKey         = "CM_TICKETMASTER_API_KEY"
	host           = "https://app.ticketmaster.com"
	eventPath      = "/discovery/v2/events"
	urlFmt         = "%s%s?classificationName=%s&dmaId=%d&localStartDateTime=%s&sort=%s&size=%v"
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

	dmaId := getDmaId(city, stateCd)

	url := fmt.Sprintf(urlFmt, host, eventPath, classification, dmaId, startDate, sort, pageSize)
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

// DMA ID is a "Designated Market Area", representing a predefined geographic region
func getDmaId(city, stateCd string) int {
	dmaMap := map[string]int{
		"atlanta-ga": 220,
	}
	key := strings.ToLower(city) + "-" + strings.ToLower(stateCd)
	return dmaMap[key]
}
