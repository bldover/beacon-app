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

type Track struct {
	Id string `json:"id"`
	Title string `json:"name"`
    Artists []Artist `json:"artists"`
}

type RankedTrack struct {
	Track Track
	Rank float64
}

type Artist struct {
    Id string `json:"id"`
	Name string `json:"name"`
}

type RankedArtist struct {
	Artist Artist
	Rank float64
}

type Client struct {
	auth *authentication
}

func NewClient() *Client {
	log.Info("Initializing Spotify client")
    client := &Client{newAuthentication()}
	log.Info("Successfully initialized Spotify client")
	return client
}

const baseUrl = "https://api.spotify.com/v1"
const limit = 50

type RequestEntity struct {
    requestUrl string
	pathParams map[string]any
}

func (c *Client) getPage(httpMethod string, reqEntity RequestEntity, response any) error {
	req, err := http.NewRequest(httpMethod, reqEntity.requestUrl, nil)
	if err != nil {
		return err
	}

	params := url.Values{}
	for name, value := range reqEntity.pathParams {
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
    Status int `json:"Status"`
	Message string `json:"Message"`
}

func (c *Client) call(req *http.Request) (*http.Response, error) {
	retries := 0
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	log.Debugf("Spotify request without auth%+v", req)

	for retries < 3 {
		authToken, err := c.auth.getAuthToken()
		req.Header.Set("Authorization", authToken)

		startTs := time.Now()
		resp, err := http.DefaultClient.Do(req)
		log.Debugf("Request response time: %v ms\n", time.Since(startTs).Milliseconds())
		if err != nil {
			return nil, err
		}

		log.Debugf("For URL %s, received response: %+v", req.URL, resp)
		if resp.StatusCode == http.StatusOK {
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
				log.Errorf("Spotify API returned high retry delay of %d, try again later", delay.Seconds())
				return nil, errors.New("exceeded rate limit and retry delay too high")
			}
			log.Debugf("Waiting %v seconds before retrying", delay.Seconds())
			time.Sleep(delay)
		default:
			log.Debug("Unexpected error, waiting 100 ms and retrying; attempt:", retries)
			time.Sleep(100 * time.Millisecond)
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

type SavedTrack struct {
    Track Track `json:"track"`
}

type SavedTrackResponse struct {
	Next string `json:"next"`
	Total int `json:"total"`
	SavedTracks []SavedTrack `json:"items"`
}
const savedTracksPath = "/me/tracks"

func (c *Client) GetSavedTracks() ([]Track, error) {
	log.Info("Request to get saved Spotify tracks")
	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	tracks := []Track{}
	savedTracksUrl := baseUrl + savedTracksPath

	// First request to get the total number of tracks
	pathParams := map[string]any{}
	pathParams["limit"] = limit
	request := RequestEntity{savedTracksUrl, pathParams}
	response := &SavedTrackResponse{}
	err := c.getPage(http.MethodGet, request, response)
	if err != nil {
		return tracks, err
	}

	savedTrackBatch := []Track{}
	for _, track := range response.SavedTracks {
		savedTrackBatch = append(savedTrackBatch, track.Track)
	}
	tracks = append(tracks, savedTrackBatch...)

	total := response.Total
	totalPages := (total + limit - 1) / limit
	log.Debugf("Spotify indicated %d total saved tracks in %d pages", total, totalPages)

	for i := 1; i < totalPages; i++ {
		wg.Add(1)
		go func(offset int) {
			defer wg.Done()
			pathParams := map[string]any{}
			pathParams["limit"] = limit
			pathParams["offset"] = offset
			request := RequestEntity{savedTracksUrl, pathParams}
			response := &SavedTrackResponse{}
			err := c.getPage(http.MethodGet, request, response)
			if err != nil {
				log.Errorf("Error fetching tracks at offset %d: %v", offset, err)
				return
			}

			savedTrackBatch := []Track{}
			for _, track := range response.SavedTracks {
				savedTrackBatch = append(savedTrackBatch, track.Track)
			}
			mu.Lock()
			tracks = append(tracks, savedTrackBatch...)
			mu.Unlock()
		}(i * limit)
		// For whatever reason, we get 500 errors only for this endpoint when sending all the
		// requests with no delay. The retry handles it, but it's faster to avoid it altogether
		time.Sleep(10 * time.Millisecond)
	}

	wg.Wait()
	retrievedCount := len(tracks)
	if retrievedCount < total {
		// tracks may still be valid for any successful batches; let the caller decide to use it or not
		errMsg := fmt.Sprintf("failed to retrieve all saved tracks, found %d/%d", retrievedCount, total)
		return tracks, errors.New(errMsg)
	}
	log.Infof("Found %v saved tracks", retrievedCount)
	return tracks, nil
}

type TopTrackResponse struct {
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

func (c *Client) GetTopTracks(timeRange TimeRange) ([]RankedTrack, error) {
	log.Info("Request to get top Spotify tracks with range:", timeRange)
	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	tracks := []RankedTrack{}
	topTracksUrl := baseUrl + topTracksPath

	// First request to get the total number of tracks
	pathParams := map[string]any{}
	pathParams["limit"] = limit
	pathParams["time_range"] = timeRange
	request := RequestEntity{topTracksUrl, pathParams}
	response := &TopTrackResponse{}
	err := c.getPage(http.MethodGet, request, response)
	if err != nil {
		return tracks, err
	}

	rankedTrackBatch := []RankedTrack{}
	for _, track := range response.TopTracks {
		rankedArtist := RankedTrack{Track: track, Rank: 0}
		rankedTrackBatch = append(rankedTrackBatch, rankedArtist)
	}
	tracks = append(tracks, rankedTrackBatch...)

	total := response.Total
	totalPages := (total + limit - 1) / limit
	log.Debugf("Spotify indicated %d total top tracks in %d pages", total, totalPages)

	for i := 1; i < totalPages; i++ {
		wg.Add(1)
		go func(offset int) {
			defer wg.Done()
			pathParams := map[string]any{}
			pathParams["limit"] = limit
			pathParams["time_range"] = timeRange
			pathParams["offset"] = offset
			request := RequestEntity{topTracksUrl, pathParams}
			response := &TopTrackResponse{}
			err := c.getPage(http.MethodGet, request, response)
			if err != nil {
				log.Errorf("Error fetching tracks at offset %d: %v", offset, err)
				return
			}

			rankedTrackBatch := []RankedTrack{}
			for _, track := range response.TopTracks {
				rankedArtist := RankedTrack{Track: track, Rank: float64(offset / total)}
				rankedTrackBatch = append(rankedTrackBatch, rankedArtist)
			}
			mu.Lock()
			tracks = append(tracks, rankedTrackBatch...)
			mu.Unlock()
		}(i * limit)
	}

	wg.Wait()
	retrievedCount := len(tracks)
	if retrievedCount < total {
		// tracks may still be valid for any successful batches; let the caller decide to use it or not
		errMsg := fmt.Sprintf("failed to retrieve all top tracks, found %d/%d", retrievedCount, total)
		return tracks, errors.New(errMsg)
	}
	log.Infof("Found %v top tracks", len(tracks))
	return tracks, nil
}

type TopArtistResponse struct {
	Next string `json:"next"`
	Total int `json:"total"`
	Offset int `json:"offset"`
	TopArtists []Artist `json:"items"`
}

const topArtistsPath = "/me/top/artists"

func (c *Client) GetTopArtists(timeRange TimeRange) ([]RankedArtist, error) {
	log.Info("Request to get top Spotify artists with range:", timeRange)
	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	artists := []RankedArtist{}
	topArtistsUrl := baseUrl + topArtistsPath

	// First request to get the total number of artists
	pathParams := map[string]any{}
	pathParams["limit"] = limit
	pathParams["time_range"] = timeRange
	request := RequestEntity{topArtistsUrl, pathParams}
	response := &TopArtistResponse{}
	err := c.getPage(http.MethodGet, request, response)
	if err != nil {
		return artists, err
	}

	rankedArtistBatch := []RankedArtist{}
	for _, artist := range response.TopArtists {
		rankedArtist := RankedArtist{Artist: artist, Rank: 0}
		rankedArtistBatch = append(rankedArtistBatch, rankedArtist)
	}
	artists = append(artists, rankedArtistBatch...)

	total := response.Total
	totalPages := (total + limit - 1) / limit
	log.Debugf("Spotify indicated %d total top artists in %d pages", total, totalPages)

	for i := 1; i < totalPages; i++ {
		wg.Add(1)
		go func(offset int) {
			defer wg.Done()
			pathParams := map[string]any{}
			pathParams["limit"] = limit
			pathParams["time_range"] = timeRange
			pathParams["offset"] = offset
			request := RequestEntity{topArtistsUrl, pathParams}
			response := &TopArtistResponse{}
			err := c.getPage(http.MethodGet, request, response)
			if err != nil {
				log.Errorf("Error fetching artists at offset %d: %v", offset, err)
				return
			}

			rankedArtistBatch := []RankedArtist{}
			for _, artist := range response.TopArtists {
				rankedArtist := RankedArtist{Artist: artist, Rank: float64(offset / total)}
				rankedArtistBatch = append(rankedArtistBatch, rankedArtist)
			}
			mu.Lock()
			artists = append(artists, rankedArtistBatch...)
			mu.Unlock()
		}(i * limit)
	}

	wg.Wait()
	retrievedCount := len(artists)
	if retrievedCount < total {
		// artists may still be valid for any successful batches; let the caller decide to use it or not
		errMsg := fmt.Sprintf("failed to retrieve all top artists, found %d/%d", retrievedCount, total)
		return artists, errors.New(errMsg)
	}
	log.Infof("Found %v top artists", len(artists))
	return artists, nil
}

type RelatedArtistResponse struct {
    Artists []Artist `json:"artists"`
}

const relatedArtistPath = "/artists/%s/related-artists"

func (c *Client) GetRelatedArtists(artist Artist) ([]Artist, error) {
	log.Info("Request to get related Spotify artists to artist:", artist)
	url := fmt.Sprintf(baseUrl + relatedArtistPath, artist.Id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.call(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var relatedArtistResp RelatedArtistResponse
	if err := json.NewDecoder(resp.Body).Decode(&relatedArtistResp); err != nil {
		errMsg := fmt.Sprintf("failed to parse response: %v", err)
		return nil, errors.New(errMsg)
	}
	log.Infof("Found related artists %v for requested artist %s", relatedArtistResp.Artists, artist.Name)
    return relatedArtistResp.Artists, nil
}

func (c *Client) GetRelatedArtistsBatch(artists []Artist) (map[Artist][]Artist, error) {
	log.Infof("Retrieved batch related artist request for %v artists", len(artists))
	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	relatedMap := map[Artist][]Artist{}
	for _, artist := range artists {
		wg.Add(1)
		go func(a Artist) {
			defer wg.Done()
			related, err := c.GetRelatedArtists(a)
			if err != nil {
				log.Errorf("Failed to retrieve related artists for %v: %v", a, err)
				return
			}
			mu.Lock()
			relatedMap[a] = related
			mu.Unlock()
		}(artist)
	}

	wg.Wait()
	successCount := len(relatedMap)
	if successCount < len(artists) {
		// artists may still be valid for any successful batches; let the caller decide to use it or not
		errMsg := fmt.Sprintf("failed to retrieve some related artists, found %d/%d", successCount, len(artists))
		return relatedMap, errors.New(errMsg)
	}
	log.Info("Finished batch related artist request")
	return relatedMap, nil
}
