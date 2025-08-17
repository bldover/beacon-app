package spotify

import (
	"concert-manager/external"
	"concert-manager/log"
	"errors"
	"fmt"
	"net/http"
)

const topTracksPath = "/me/top/tracks"

func (c *Client) TopTracks(timeRange external.TimeRange) ([]external.Track, error) {
	log.Info("Request to get top Spotify tracks with range:", timeRange)
	tracks := []track{}
	topTracksUrl := baseUrl + topTracksPath

	// First request to get the total number of tracks
	queryParams := map[string]any{}
	queryParams["limit"] = limit
	queryParams["time_range"] = timeRange
	request := RequestEntity{topTracksUrl, queryParams}
	response := &topTrackResponse{}
	err := c.call(http.MethodGet, request, response)
	if err != nil {
		return mapSpotifyTracks(tracks), err
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
		err := c.call(http.MethodGet, request, response)
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
		return mapSpotifyTracks(tracks), errors.New(errMsg)
	}
	log.Infof("Found %v top tracks", len(tracks))
	return mapSpotifyTracks(tracks), nil
}

const topArtistsPath = "/me/top/artists"

func (c *Client) TopArtists(timeRange external.TimeRange) ([]external.Artist, error) {
	log.Info("Request to get top Spotify artists with range:", timeRange)
	artists := []artist{}
	topArtistsUrl := baseUrl + topArtistsPath

	// First request to get the total number of artists
	queryParams := map[string]any{}
	queryParams["limit"] = limit
	queryParams["time_range"] = timeRange
	request := RequestEntity{topArtistsUrl, queryParams}
	response := &topArtistResponse{}
	err := c.call(http.MethodGet, request, response)
	if err != nil {
		return mapSpotifyArtists(artists), err
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
		err := c.call(http.MethodGet, request, response)
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
		return mapSpotifyArtists(artists), errors.New(errMsg)
	}
	log.Infof("Found %v top artists", len(artists))
	return mapSpotifyArtists(artists), nil
}
