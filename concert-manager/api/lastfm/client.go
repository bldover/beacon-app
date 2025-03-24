package lastfm

import (
	"concert-manager/api"
	"concert-manager/log"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const baseUrl = "http://ws.audioscrobbler.com/2.0/"
const apiKeyEnv = "CM_LASTFM_API_KEY"

type Client struct {
	apiKey string
}

func NewClient() *Client {
	apiKey := os.Getenv(apiKeyEnv)
	if apiKey == "" {
		log.Fatalf("%s env var must be set", apiKeyEnv)
	}
	return &Client{
		apiKey: apiKey,
	}
}

func toArtist(artist artist) api.Artist {
	rank, err := strconv.ParseFloat(artist.Rank, 64)
	if err != nil {
		log.Errorf("invalid match value for lastfm similar artist: %s", artist)
		rank = 0
	}
    return api.Artist{
		Name: artist.Name,
		Rank: float64(rank),
	}
}

const maxRetryCount = 3

type requestEntity struct {
    queryParams map[string]any
}

type errorResponse struct {
    Message string `json:"message"`
	Code int `json:"error"`
}

func (c *Client) call(reqEntity requestEntity, response any) error {
	req, err := http.NewRequest(http.MethodGet, baseUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("autocorrect", "1")

	queryParams := reqEntity.queryParams
	queryParams["api_key"] = c.apiKey
	queryParams["format"] = "json"
	params := url.Values{}
	for name, value := range reqEntity.queryParams {
		params.Set(name, fmt.Sprintf("%v", value))
	}
	req.URL.RawQuery = params.Encode()

	retries := 0
	for retries < maxRetryCount {
		startTs := time.Now()
		resp, err := http.DefaultClient.Do(req)
		log.Debugf("Request response time: %v ms\n", time.Since(startTs).Milliseconds())
		if err != nil {
			return err
		}

		log.Debugf("For URL %s, received response: %+v", req.URL, resp)
		if resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
				errMsg := fmt.Sprintf("failed to parse response: %v", err)
				return errors.New(errMsg)
			}
			return nil
		}

		retries += 1
		errorResp := &errorResponse{}
		if err := json.NewDecoder(resp.Body).Decode(errorResp); err != nil {
			log.Error("Failed to decode error response", resp)
		} else {
			log.Error("Received LastFM error response", errorResp)
		}

		if resp.StatusCode == http.StatusTooManyRequests && retries < maxRetryCount {
			delay := 1 * time.Second
			time.Sleep(delay)
		}
	}
	return errors.New("max retries exceeded calling LastFM URL: " + req.URL.Host + req.URL.Path + "?" + req.URL.RawQuery)
}

type relatedArtistResponse struct {
	Similar similarArtistList `json:"similarartists"`
}

type similarArtistList struct {
	Artists []artist `json:"artist"`
}

type artist struct {
    Name string `json:"name"`
	Rank string `json:"match"`
}

func (c *Client) GetRelatedArtists(artists []api.Artist) (map[api.Artist][]api.Artist, error) {
	related := map[api.Artist][]api.Artist{}
	successCount := 0

	for _, artist := range artists {
		queryParams := map[string]any{}
		queryParams["method"] = "artist.getsimilar"
		queryParams["artist"] = artist.Name

		request := requestEntity{queryParams}
		response := &relatedArtistResponse{}
		err := c.call(request, response)
		if err != nil {
			log.Errorf("Failed to retrieve related artists: %s", err)
			continue
		}

		if response.Similar.Artists == nil || len(response.Similar.Artists) == 0 {
			log.Errorf("No related artists found for %s", artist.Name)
			continue
		}

		related[artist] = []api.Artist{}
		for _, relatedArtist := range response.Similar.Artists {
			related[artist] = append(related[artist], toArtist(relatedArtist))
		}
		successCount += 1
	}
	if successCount < len(artists) {
		errMsg := fmt.Sprintf("failed to retrieve some related artists, found %d/%d", successCount, len(artists))
		return related, errors.New(errMsg)
	}

	return related, nil
}
