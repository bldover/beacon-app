package lastfm

import (
	"concert-manager/log"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

const baseUrl = "http://ws.audioscrobbler.com/2.0/"
const apiKeyEnv = "CM_LASTFM_API_KEY"

const userAgent = "Beacon/2.0"

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

const maxRetryCount = 3

type requestEntity struct {
	queryParams map[string]any
}

type errorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"error"`
}

func (c *Client) call(reqEntity requestEntity, response any) error {
	req, err := http.NewRequest(http.MethodGet, baseUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("autocorrect", "1")
	req.Header.Set("User-Agent", "Beacon/1.0")

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
		log.Debugf("Request response time: %v ms", time.Since(startTs).Milliseconds())
		if err != nil {
			return err
		}

		log.Debugf("For URL %v, received response: %+v", stripApiKey(*req.URL), resp)
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

func stripApiKey(reqUrl url.URL) url.URL {
	values, _ := url.ParseQuery(reqUrl.RawQuery)
	values.Del("api_key")
	reqUrl.RawQuery = values.Encode()
	return reqUrl
}
