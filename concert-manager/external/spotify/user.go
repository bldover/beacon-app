package spotify

import (
	"concert-manager/external"
	"concert-manager/log"
	"errors"
	"fmt"
	"net/http"
)

const savedTracksPath = "/me/tracks"

func (c *Client) SavedTracks() ([]external.Track, error) {
	log.Info("Request to get saved Spotify tracks")
	tracks := []track{}
	savedTracksUrl := baseUrl + savedTracksPath

	// First request to get the total number of tracks
	queryParams := map[string]any{}
	queryParams["limit"] = limit
	request := RequestEntity{savedTracksUrl, queryParams}
	response := &savedTrackResponse{}
	err := c.call(http.MethodGet, request, response)
	if err != nil {
		return mapSpotifyTracks(tracks), err
	}

	savedTrackBatch := []track{}
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
		err := c.call(http.MethodGet, request, response)
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
		return mapSpotifyTracks(tracks), errors.New(errMsg)
	}
	log.Infof("Found %v saved tracks", retrievedCount)

	return mapSpotifyTracks(tracks), nil
}
