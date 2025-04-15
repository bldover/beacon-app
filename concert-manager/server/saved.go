package server

import (
	"concert-manager/data"
	"concert-manager/log"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func (s *Server) handleVenues(w http.ResponseWriter, r *http.Request) (any, int, error) {
	switch r.Method {
    case http.MethodGet:
		venues := s.VenueCache.GetVenues()
		return venues, 0, nil
	case http.MethodPost:
		var venue data.Venue
		if err := json.NewDecoder(r.Body).Decode(&venue); err != nil {
			return nil, http.StatusBadRequest, errors.New("invalid body")
		}
		savedVenue, err := s.VenueCache.AddVenue(venue)
		if err != nil {
			errMsg := fmt.Sprintf("failed to save venue: %v", err)
			return nil, http.StatusInternalServerError, errors.New(errMsg)
		}
		return savedVenue, 0, nil
	case http.MethodPut:
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) != 4 {
			return nil, http.StatusBadRequest, errors.New("missing venue ID in path")
		}
		id := pathParts[3]
		if len(id) == 0 {
			return nil, http.StatusBadRequest, errors.New("missing venue ID in path")
		}
		var venue data.Venue
		if err := json.NewDecoder(r.Body).Decode(&venue); err != nil {
			return nil, http.StatusBadRequest, errors.New("invalid body")
		}
		err := s.VenueCache.UpdateVenue(id, venue)
		if err != nil {
			errMsg := fmt.Sprintf("failed to update venue: %v", err)
			return nil, http.StatusInternalServerError, errors.New(errMsg)
		}
		return nil, 0, nil
	case http.MethodDelete:
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) != 4 {
			return nil, http.StatusBadRequest, errors.New("missing venue ID in path")
		}
		id := pathParts[3]
		if len(id) == 0 {
			return nil, http.StatusBadRequest, errors.New("missing venue ID in path")
		}
		if err := s.VenueCache.DeleteVenue(id); err != nil {
			errMsg := fmt.Sprintf("failed to delete venue: %v", err)
			return nil, http.StatusInternalServerError, errors.New(errMsg)
		}
		return nil, 0, nil
	}
	return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
}

func (s *Server) handleArtists(w http.ResponseWriter, r *http.Request) (any, int, error) {
	switch r.Method {
    case http.MethodGet:
		artists := s.ArtistCache.GetArtists()
		return artists, 0, nil
	case http.MethodPost:
		var artist data.Artist
		if err := json.NewDecoder(r.Body).Decode(&artist); err != nil {
			return nil, http.StatusBadRequest, errors.New("invalid body")
		}
		savedArtist, err := s.ArtistCache.AddArtist(artist)
		if err != nil {
			errMsg := fmt.Sprintf("failed to save artist: %v", err)
			return nil, http.StatusInternalServerError, errors.New(errMsg)
		}
		return savedArtist, 0, nil
	case http.MethodPut:
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) != 4 {
			return nil, http.StatusBadRequest, errors.New("missing artist ID in path")
		}
		id := pathParts[3]
		if len(id) == 0 {
			return nil, http.StatusBadRequest, errors.New("missing artist ID in path")
		}
		var artist data.Artist
		if err := json.NewDecoder(r.Body).Decode(&artist); err != nil {
			return nil, http.StatusBadRequest, errors.New("invalid body")
		}
		err := s.ArtistCache.UpdateArtist(id, artist)
		if err != nil {
			errMsg := fmt.Sprintf("failed to update artist: %v", err)
			return nil, http.StatusInternalServerError, errors.New(errMsg)
		}
		return nil, 0, nil
	case http.MethodDelete:
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) != 4 {
			return nil, http.StatusBadRequest, errors.New("missing artist ID in path")
		}
		id := pathParts[3]
		if len(id) == 0 {
			return nil, http.StatusBadRequest, errors.New("missing artist ID in path")
		}
		if err := s.ArtistCache.DeleteArtist(id); err != nil {
			errMsg := fmt.Sprintf("failed to delete artist: %v", err)
			return nil, http.StatusInternalServerError, errors.New(errMsg)
		}
		return nil, 0, nil
	}
	return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
}

func (s *Server) handleSavedEvents(w http.ResponseWriter, r *http.Request) (any, int, error) {
	switch r.Method {
    case http.MethodGet:
		events := s.SavedEventCache.GetSavedEvents()
		id := r.URL.Query().Get("id")
		if id == "" {
			return events, 0, nil
		}
		for _, event := range events {
			if event.Id == id {
				return []data.Event{event}, 0, nil
			}
		}
		errMsg := fmt.Sprintf("event with ID %s not found", id)
		return nil, http.StatusNotFound, errors.New(errMsg)
	case http.MethodPost:
		var event data.Event
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			return nil, http.StatusBadRequest, errors.New("invalid body")
		}
		savedEvent, err := s.SavedEventCache.AddSavedEvent(event)
		if err != nil {
			errMsg := fmt.Sprintf("failed to save event: %v", err)
			return nil, http.StatusInternalServerError, errors.New(errMsg)
		}
		return savedEvent, 0, nil
	case http.MethodDelete:
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) != 5 {
			return nil, http.StatusBadRequest, errors.New("missing event ID in path")
		}
		id := pathParts[4]
		if len(id) == 0 {
			return nil, http.StatusBadRequest, errors.New("missing event ID in path")
		}
		if err := s.SavedEventCache.DeleteSavedEvent(id); err != nil {
			errMsg := fmt.Sprintf("failed to delete event: %v", err)
			return nil, http.StatusInternalServerError, errors.New(errMsg)
		}
		return nil, 0, nil
	}
	return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
}

func (s *Server) refreshSavedEvents(w http.ResponseWriter, r *http.Request) (any, int, error) {
    if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
	}

	err := s.SavedEventCache.RefreshSavedEvents()
	if err != nil {
		log.Errorf("Failed to refresh saved events %v", err)
		return nil, http.StatusInternalServerError, errors.New("failed to refresh saved event cache")
	}
	return nil, 0, nil
}

func (s *Server) refreshArtists(w http.ResponseWriter, r *http.Request) (any, int, error) {
    if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
	}

	err := s.ArtistCache.RefreshArtists()
	if err != nil {
		log.Errorf("Failed to refresh artists %v", err)
		return nil, http.StatusInternalServerError, errors.New("failed to refresh artists cache")
	}
	return nil, 0, nil
}

func (s *Server) refreshVenues(w http.ResponseWriter, r *http.Request) (any, int, error) {
    if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
	}

	err := s.VenueCache.RefreshVenues()
	if err != nil {
		log.Errorf("Failed to refresh venues %v", err)
		return nil, http.StatusInternalServerError, errors.New("failed to refresh venues cache")
	}
	return nil, 0, nil
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		errMsg := fmt.Sprintf("unable to parse request file: %v", err)
		return nil, http.StatusBadRequest, errors.New(errMsg)
	}

    rows, err := s.EventLoader.Upload(r.Context(), file)
	if err != nil {
		errMsg := fmt.Sprintf("error occurred during upload processing: %v", err)
		return nil, http.StatusBadRequest, errors.New(errMsg)
	}
	return fmt.Sprintf("Successfully uploaded %d rows", rows), 0, nil
}

func (s *Server) loadLastFmGenres(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodGet {
		return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
	}

	queryParams := r.URL.Query()
	artists := queryParams["artists"]

	updateCount, err := s.ArtistInfoLoader.PopulateLastFmData(r.Context(), artists)
	if err != nil {
		errMsg := fmt.Sprintf("some artist genres failed to refresh, %s", err)
		return nil, http.StatusInternalServerError, errors.New(errMsg)
	}
	return fmt.Sprintf("Successfully refreshed genres for %d artists", updateCount), 0, nil
}

func (s *Server) loadUserGenres(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodGet {
		return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
	}

	updateCount, err := s.ArtistInfoLoader.PopulateUserGenres(r.Context())
	if err != nil {
		errMsg := fmt.Sprintf("some artist user genres failed to populate, %v", err)
		return nil, http.StatusInternalServerError, errors.New(errMsg)
	}
	return fmt.Sprintf("Successfully loaded user genres for %d artists", updateCount), 0, nil
}
