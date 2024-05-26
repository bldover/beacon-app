package spotify

import (
	"concert-manager/log"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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
    return &Client{newAuthentication()}
}

const baseUrl = "https://api.spotify.com/v1"
const limit = 50

type SavedTrackResponse struct {
	Next string `json:"next"`
	Total int `json:"total"`
	SavedTracks []struct {
		Track Track `json:"track"`
	} `json:"items"`
}

const savedTracksPath = "/me/tracks"

func (s *Client) GetSavedTracks() ([]Track, error) {
	log.Info("Request to get saved Spotify tracks")
	var tracks []Track
	savedTrackUrl := baseUrl + savedTracksPath
	for savedTrackUrl != "" {
		req, err := http.NewRequest(http.MethodGet, savedTrackUrl, nil)
		if err != nil {
			return tracks, err
		}
		if req.URL.RawQuery == "" {
			params := url.Values{}
			params.Set("limit", strconv.Itoa(limit))
			req.URL.RawQuery = params.Encode()
		}

		resp, err := s.call(req)
		if err != nil {
			return tracks, err
		}
		defer resp.Body.Close()

		var trackResponse SavedTrackResponse
		if err := json.NewDecoder(resp.Body).Decode(&trackResponse); err != nil {
			errMsg := fmt.Sprintf("failed to parse response: %v", err)
			return tracks, errors.New(errMsg)
		}

		if tracks == nil {
			tracks = make([]Track, 0, trackResponse.Total)
		}
		for _, savedTrack := range trackResponse.SavedTracks {
			tracks = append(tracks, savedTrack.Track)
		}

		savedTrackUrl = trackResponse.Next
	}
	log.Infof("Found %v saved tracks", len(tracks))
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

func (s *Client) GetTopTracks(timeRange TimeRange) ([]RankedTrack, error) {
	log.Info("Request to get top Spotify tracks with range:", timeRange)
	var tracks []RankedTrack
	topTracksUrl := baseUrl + topTracksPath
	for topTracksUrl != "" {
		req, err := http.NewRequest(http.MethodGet, topTracksUrl, nil)
		if err != nil {
			return tracks, err
		}
		if req.URL.RawQuery == "" {
			params := url.Values{}
			params.Set("limit", strconv.Itoa(limit))
			params.Set("time_range", string(timeRange))
			req.URL.RawQuery = params.Encode()
		}

		resp, err := s.call(req)
		if err != nil {
			return tracks, err
		}
		defer resp.Body.Close()

		var trackResponse TopTrackResponse
		if err := json.NewDecoder(resp.Body).Decode(&trackResponse); err != nil {
			errMsg := fmt.Sprintf("failed to parse response: %v", err)
			return tracks, errors.New(errMsg)
		}

		if tracks == nil {
			tracks = make([]RankedTrack, 0, trackResponse.Total)
		}
		for _, topTrack := range trackResponse.TopTracks {
			rankedTrack := RankedTrack{
				Track: topTrack,
				Rank: float64(trackResponse.Offset / limit),
			}
			tracks = append(tracks, rankedTrack)
		}

		topTracksUrl = trackResponse.Next
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

func (s *Client) GetTopArtists(timeRange TimeRange) ([]RankedArtist, error) {
	log.Info("Request to get top Spotify artists with range:", timeRange)
	var artists []RankedArtist
	topArtistsUrl := baseUrl + topArtistsPath
	for topArtistsUrl != "" {
		req, err := http.NewRequest(http.MethodGet, topArtistsUrl, nil)
		if err != nil {
			return artists, err
		}
		if req.URL.RawQuery == "" {
			params := url.Values{}
			params.Set("limit", strconv.Itoa(limit))
			params.Set("time_range", string(timeRange))
			req.URL.RawQuery = params.Encode()
		}

		resp, err := s.call(req)
		if err != nil {
			return artists, err
		}
		defer resp.Body.Close()

		var artistResponse TopArtistResponse
		if err := json.NewDecoder(resp.Body).Decode(&artistResponse); err != nil {
			errMsg := fmt.Sprintf("failed to parse response: %v", err)
			return artists, errors.New(errMsg)
		}

		if artists == nil {
			artists = make([]RankedArtist, 0, artistResponse.Total)
		}
		for _, topArtist := range artistResponse.TopArtists {
			rankedArtist := RankedArtist{
				Artist: topArtist,
				Rank: float64(artistResponse.Offset / limit),
			}
			artists = append(artists, rankedArtist)
		}

		topArtistsUrl = artistResponse.Next
	}
	log.Infof("Found %v top artists", len(artists))
	return artists, nil
}

type RelatedArtistResponse struct {
    Artists []Artist `json:"artists"`
}

const relatedArtistPath = "/artists/%s/related-artists"

func (s *Client) GetRelatedArtists(artist Artist) ([]Artist, error) {
	log.Infof("Request to get related Spotify artists to %v", artist)
	url := fmt.Sprintf(baseUrl + relatedArtistPath, artist.Id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.call(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var relatedArtistResp RelatedArtistResponse
	if err := json.NewDecoder(resp.Body).Decode(&relatedArtistResp); err != nil {
		errMsg := fmt.Sprintf("failed to parse response: %v", err)
		return nil, errors.New(errMsg)
	}
    return relatedArtistResp.Artists, nil
}

type errorResponse struct {
    Status int `json:"status"`
	Message string `json:"message"`
}

func (s *Client) call(req *http.Request) (*http.Response, error) {
	retries := 0
	if s.auth.accessToken == "" {
		if err := s.auth.refresh(); err != nil {
			errMsg := fmt.Sprintf("failed to refresh Spotify token: %v", err)
			return nil, errors.New(errMsg)
		}
	}
	req.Header.Set("Authorization", s.auth.accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	log.Debugf("%+v", req)
	for retries < 3 {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		log.Debug(resp)
		if resp.StatusCode == http.StatusOK {
			return resp, nil
		}

		var errorResp errorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			log.Error("Failed to decode error response")
		} else {
			log.Error("Received Spotify error response", errorResp)
		}

		switch resp.StatusCode {
		case http.StatusUnauthorized:
			if err := s.auth.refresh(); err != nil {
				errMsg := fmt.Sprintf("failed to refresh Spotify token: %v", err)
				return nil, errors.New(errMsg)
			}
			req.Header.Set("Authentication", s.auth.accessToken)
		case http.StatusTooManyRequests:
			delay := getDelay(resp)
			log.Debugf("Waiting %v seconds before retrying", delay)
			time.Sleep(time.Duration(delay))
		default:
			return nil, errors.New("unexpected error")
		}
		retries++
	}
	return nil, errors.New("max retries exceeded calling Spotify URL: " + req.URL.Host + req.URL.Path)
}

func getDelay(resp *http.Response) int {
	delayHeader := resp.Header.Get("Retry-After")
	delay, err := strconv.Atoi(delayHeader)
	if err != nil {
		log.Error("Failed to parse Retry-After header for Spotify 429 response", delayHeader)
		delay = 30
	}
	return delay
}
