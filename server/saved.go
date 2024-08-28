package server

import (
	"concert-manager/data"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			return nil, http.StatusMethodNotAllowed, errors.New("missing id request parameter")
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
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			return nil, http.StatusBadRequest, errors.New("missing id request parameter")
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
		id := r.URL.Query().Get("id")
		if id == "" {
			return nil, http.StatusBadRequest, errors.New("missing id request parameter")
		}
		if err := s.SavedEventCache.DeleteSavedEvent(id); err != nil {
			errMsg := fmt.Sprintf("Failed to delete event: %v", err)
			return nil, http.StatusInternalServerError, errors.New(errMsg)
		}
		return nil, 0, nil
	}
	return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
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

    rows, err := s.Loader.Upload(r.Context(), file)
	if err != nil {
		errMsg := fmt.Sprintf("error occurred during upload processing: %v", err)
		return nil, http.StatusBadRequest, errors.New(errMsg)
	}
	return fmt.Sprintf("Successfully uploaded %d rows", rows), 0, nil
}
