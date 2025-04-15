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
	auth *authentication
	threadInitDelay time.Duration
	retryStrategy RetryStrategy
}

const defaultBackoffIncrement = 50 * time.Millisecond

type RetryStrategy struct {
	delay time.Duration
	backoff time.Duration
	increment time.Duration
	lock *sync.Mutex
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
    requestUrl string
	queryParams map[string]any
}

func (c *Client) getPage(httpMethod string, reqEntity RequestEntity, response any) error {
	req, err := http.NewRequest(httpMethod, reqEntity.requestUrl, nil)
	if err != nil {
		return err
	}

	params := url.Values{}
	for name, value := range reqEntity.queryParams {
		params.Set(name, fmt.Sprintf("%v", value))
	}
	req.URL.RawQuery = params.Encode()

	resp, err := c.call(req)
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
	Error struct {
		Status int `json:"Status"`
		Message string `json:"Message"`
	}
}

func (c *Client) call(req *http.Request) (*http.Response, error) {
	retries := 0
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	log.Debugf("Spotify request without auth%+v", req)

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
			return nil, err
		}

		log.Debugf("For URL %s, received response: %+v", req.URL, resp)
		if resp.StatusCode == http.StatusOK {
			c.retryStrategy.lock.Lock()
			c.retryStrategy.backoff -= c.retryStrategy.increment * time.Duration(retries + 1)
			log.Debugf("Success request, new backoff %f seconds", c.retryStrategy.backoff.Seconds())
			c.retryStrategy.lock.Unlock()
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
				return nil, errors.New("exceeded rate limit and retry delay too high")
			}

			c.retryStrategy.lock.Lock()
			c.retryStrategy.delay = delay
			c.retryStrategy.lock.Unlock()
			log.Debugf("TooManyRequests; delay=%f, backoff=%f, attempt: %d", delay.Seconds(), c.retryStrategy.backoff.Seconds(), retries)
		default:
			log.Debug("Unexpected error; attempt:", retries)
		}
		retries++
	}
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

type savedTrackResponse struct {
	Next string `json:"next"`
	Total int `json:"total"`
	SavedTracks []trackInfo `json:"items"`
}

type trackInfo struct {
    Track Track `json:"track"`
}

type Track struct {
	Id string `json:"id"`
	Title string `json:"name"`
    Artists []Artist `json:"artists"`
}

type Artist struct {
    Id string `json:"id"`
	Name string `json:"name"`
}

const savedTracksPath = "/me/tracks"

func (c *Client) GetSavedTracks() ([]Track, error) {
	log.Info("Request to get saved Spotify tracks")
	log.Debug("TEST_MODE:", TEST_MODE)
	tracks := []Track{}
	savedTracksUrl := baseUrl + savedTracksPath

	// First request to get the total number of tracks
	queryParams := map[string]any{}
	queryParams["limit"] = limit
	request := RequestEntity{savedTracksUrl, queryParams}
	response := &savedTrackResponse{}
	err := c.getPage(http.MethodGet, request, response)
	if err != nil {
		return tracks, err
	}

	savedTrackBatch := []Track{}
	for _, savedTrack := range response.SavedTracks {
		savedTrackBatch = append(savedTrackBatch, savedTrack.Track)
	}
	tracks = append(tracks, savedTrackBatch...)

	total := response.Total
	totalPages := (total + limit - 1) / limit
	log.Debugf("Spotify indicated %d total saved tracks in %d pages", total, totalPages)

	for i := 1; i < totalPages; i++ {
		offset := i * limit
		queryParams := map[string]any{}
		queryParams["limit"] = limit
		queryParams["offset"] = offset
		request := RequestEntity{savedTracksUrl, queryParams}
		response := &savedTrackResponse{}
		err := c.getPage(http.MethodGet, request, response)
		if err != nil {
			log.Errorf("Error fetching tracks at offset %d: %v", offset, err)
			continue
		}

		for _, savedTrack := range response.SavedTracks {
			tracks = append(tracks, savedTrack.Track)
		}

		if TEST_MODE && i >= testModePageLimit {
			break
		}
	}

	retrievedCount := len(tracks)
	if retrievedCount < total && !TEST_MODE {
		// tracks may still be valid for any successful batches; let the caller decide to use it or not
		errMsg := fmt.Sprintf("failed to retrieve all saved tracks, found %d/%d", retrievedCount, total)
		return tracks, errors.New(errMsg)
	}
	log.Infof("Found %v saved tracks", retrievedCount)

	return tracks, nil
}

type topTrackResponse struct {
	Next string `json:"next"`
	Total int `json:"total"`
	Offset int `json:"offset"`
	TopTracks []Track `json:"items"`
}

const topTracksPath = "/me/top/tracks"

type TimeRange string
const LongTerm = "long_term"
const MediumTerm = "medium_term"
const ShortTerm = "short_term"

func (c *Client) GetTopTracks(timeRange TimeRange) ([]Track, error) {
	log.Info("Request to get top Spotify tracks with range:", timeRange)
	tracks := []Track{}
	topTracksUrl := baseUrl + topTracksPath

	// First request to get the total number of tracks
	queryParams := map[string]any{}
	queryParams["limit"] = limit
	queryParams["time_range"] = timeRange
	request := RequestEntity{topTracksUrl, queryParams}
	response := &topTrackResponse{}
	err := c.getPage(http.MethodGet, request, response)
	if err != nil {
		return tracks, err
	}

	tracks = append(tracks, response.TopTracks...)

	total := response.Total
	totalPages := (total + limit - 1) / limit
	log.Debugf("Spotify indicated %d total top tracks in %d pages", total, totalPages)

	for i := 1; i < totalPages; i++ {
		offset := i * limit
		queryParams := map[string]any{}
		queryParams["limit"] = limit
		queryParams["time_range"] = timeRange
		queryParams["offset"] = offset
		request := RequestEntity{topTracksUrl, queryParams}
		response := &topTrackResponse{}
		err := c.getPage(http.MethodGet, request, response)
		if err != nil {
			log.Errorf("Error fetching tracks at offset %d: %v", offset, err)
			continue
		}

		tracks = append(tracks, response.TopTracks...)

		if TEST_MODE && i >= testModePageLimit {
			break
		}
	}

	retrievedCount := len(tracks)
	if retrievedCount < total && !TEST_MODE {
		// tracks may still be valid for any successful batches; let the caller decide to use it or not
		errMsg := fmt.Sprintf("failed to retrieve all top tracks, found %d/%d", retrievedCount, total)
		return tracks, errors.New(errMsg)
	}
	log.Infof("Found %v top tracks", len(tracks))
	return tracks, nil
}

type topArtistResponse struct {
	Next string `json:"next"`
	Total int `json:"total"`
	Offset int `json:"offset"`
	TopArtists []Artist `json:"items"`
}

const topArtistsPath = "/me/top/artists"

func (c *Client) GetTopArtists(timeRange TimeRange) ([]Artist, error) {
	log.Info("Request to get top Spotify artists with range:", timeRange)
	artists := []Artist{}
	topArtistsUrl := baseUrl + topArtistsPath

	// First request to get the total number of artists
	queryParams := map[string]any{}
	queryParams["limit"] = limit
	queryParams["time_range"] = timeRange
	request := RequestEntity{topArtistsUrl, queryParams}
	response := &topArtistResponse{}
	err := c.getPage(http.MethodGet, request, response)
	if err != nil {
		return artists, err
	}

	artists = append(artists, response.TopArtists...)

	total := response.Total
	totalPages := (total + limit - 1) / limit
	log.Debugf("Spotify indicated %d total top artists in %d pages", total, totalPages)

	for i := 1; i < totalPages; i++ {
		offset := i * limit
		queryParams := map[string]any{}
		queryParams["limit"] = limit
		queryParams["time_range"] = timeRange
		queryParams["offset"] = offset
		request := RequestEntity{topArtistsUrl, queryParams}
		response := &topArtistResponse{}
		err := c.getPage(http.MethodGet, request, response)
		if err != nil {
			log.Errorf("Error fetching artists at offset %d: %v", offset, err)
			continue
		}

		artists = append(artists, response.TopArtists...)

		if TEST_MODE && i >= testModePageLimit {
			break
		}
	}

	retrievedCount := len(artists)
	if retrievedCount < total && !TEST_MODE {
		// artists may still be valid for any successful batches; let the caller decide to use it or not
		errMsg := fmt.Sprintf("failed to retrieve all top artists, found %d/%d", retrievedCount, total)
		return artists, errors.New(errMsg)
	}
	log.Infof("Found %v top artists", len(artists))
	return artists, nil
}
