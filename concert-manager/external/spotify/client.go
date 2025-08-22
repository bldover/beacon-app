package spotify

import (
	"concert-manager/log"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

var TEST_MODE = false

const testModePageLimit = 20

type Client struct {
	auth            *authentication
	threadInitDelay time.Duration
	retryStrategy   RetryStrategy
}

const defaultBackoffIncrement = 50 * time.Millisecond

type RetryStrategy struct {
	delay     time.Duration
	backoff   time.Duration
	increment time.Duration
	lock      *sync.Mutex
}

func NewClient() *Client {
	log.Info("Initializing Spotify client")
	retryStrategy := RetryStrategy{0, 0, defaultBackoffIncrement, &sync.Mutex{}}
	client := &Client{newAuthentication(), defaultThreadInitDelay, retryStrategy}
	log.Info("Successfully initialized Spotify client")
	return client
}

const baseUrl = "https://api.spotify.com/v1"
const limit = 50
const defaultThreadInitDelay = 0 * time.Millisecond

type RequestEntity struct {
	requestUrl  string
	queryParams map[string]any
}

func (c *Client) call(httpMethod string, reqEntity RequestEntity, response any) error {
	req, err := http.NewRequest(httpMethod, reqEntity.requestUrl, nil)
	if err != nil {
		return err
	}

	params := url.Values{}
	for name, value := range reqEntity.queryParams {
		params.Set(name, fmt.Sprintf("%v", value))
	}
	req.URL.RawQuery = params.Encode()

	resp, err := c.execute(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		errMsg := fmt.Sprintf("failed to parse response: %v", err)
		return errors.New(errMsg)
	}

	return nil
}

type errorResponse struct {
	ErrorDetails struct {
		Status  int    `json:"Status"`
		Message string `json:"Message"`
	} `json:"error"`
}

func (e errorResponse) Error() string {
	return fmt.Sprintf("error from Spotify with code: %d, message: %s", e.ErrorDetails.Status, e.ErrorDetails.Message)
}

func (c *Client) execute(req *http.Request) (*http.Response, error) {
	retries := 0
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	log.Debugf("Calling Spotify URL: %s", req.URL)

	for retries < 3 {
		c.retryStrategy.lock.Lock()
		delay := c.retryStrategy.delay + c.retryStrategy.backoff
		c.retryStrategy.backoff += c.retryStrategy.increment
		c.retryStrategy.lock.Unlock()
		if delay > 0 {
			log.Debugf("Waiting %v seconds before request", delay.Seconds())
			time.Sleep(delay)
		}

		authToken, err := c.auth.getAuthToken()
		if err != nil {
			log.Errorf("unable to retrieve auth token: %v", err)
			retries++
			continue
		}

		req.Header.Set("Authorization", authToken)
		startTs := time.Now()
		resp, err := http.DefaultClient.Do(req)
		log.Debugf("Request response time: %v ms\n", time.Since(startTs).Milliseconds())
		if err != nil {
			c.clearBackoff(retries, true)
			return nil, err
		}

		log.Debugf("For URL %s, received response: %+v", req.URL, resp)
		if resp.StatusCode == http.StatusOK {
			c.clearBackoff(retries, true)
			return resp, nil
		}

		errorResp := &errorResponse{}
		if err := json.NewDecoder(resp.Body).Decode(errorResp); err != nil {
			log.Error("Failed to decode error response", resp)
		} else {
			log.Error("Received Spotify error response", errorResp)
		}

		switch resp.StatusCode {
		case http.StatusUnauthorized:
			c.auth.markAuthExpired(authToken)
		case http.StatusTooManyRequests:
			delay := getDelay(resp)
			if delay > (30 * time.Second) {
				log.Errorf("Spotify API returned high retry delay of %v, try again later", delay.Seconds())
				c.clearBackoff(retries, false)
				return nil, errors.New("exceeded rate limit and retry delay too high")
			}

			c.retryStrategy.lock.Lock()
			c.retryStrategy.delay = delay
			c.retryStrategy.lock.Unlock()
			log.Debugf("TooManyRequests; delay=%f, backoff=%f, attempt: %d", delay.Seconds(), c.retryStrategy.backoff.Seconds(), retries)
		default:
			code := resp.StatusCode
			if code < 500 {
				c.clearBackoff(retries, false)
				log.Errorf("non-retryable error code %d calling Spotify URL: %s", code, req.URL.Host+req.URL.Path)
				return nil, errorResp
			}
			log.Debug("Unexpected Spotify server error; attempt:", retries)
		}
		retries++
	}
	c.clearBackoff(retries-1, false)
	return nil, errors.New("max retries exceeded calling Spotify URL: " + req.URL.Host + req.URL.Path)
}

func getDelay(resp *http.Response) time.Duration {
	delayHeader := resp.Header.Get("Retry-After")
	delay, err := strconv.Atoi(delayHeader)
	if err != nil {
		log.Error("Failed to parse Retry-After header for Spotify 429 response", delayHeader)
		delay = 30
	}
	return (time.Duration(delay) + 1) * time.Second
}

func (c *Client) clearBackoff(retries int, success bool) {
	c.retryStrategy.lock.Lock()
	requestBackoff := c.retryStrategy.increment * time.Duration(retries+1)
	c.retryStrategy.backoff -= requestBackoff
	log.Debugf("new backoff %f seconds", c.retryStrategy.backoff.Seconds())
	if success {
		c.retryStrategy.delay = 0
	}
	c.retryStrategy.lock.Unlock()
}
