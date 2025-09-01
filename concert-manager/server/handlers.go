package server

import (
	"concert-manager/domain"
	"concert-manager/log"
	"concert-manager/ranker"
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
		var venue domain.Venue
		if err := json.NewDecoder(r.Body).Decode(&venue); err != nil {
			return nil, http.StatusBadRequest, errors.New("invalid body")
		}
		savedVenue, err := s.VenueCache.AddVenue(venue)
		if err != nil {
			errMsg := fmt.Sprintf("failed to save venue: %v", err)
			return nil, http.StatusInternalServerError, errors.New(errMsg)
		}
		err = s.SyncService.SyncVenueAdd(savedVenue.ID.Primary)
		if err != nil {
			errMsg := fmt.Sprintf("failed to sync change for venue: %v", err)
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
		var venue domain.Venue
		if err := json.NewDecoder(r.Body).Decode(&venue); err != nil {
			return nil, http.StatusBadRequest, errors.New("invalid body")
		}
		err := s.VenueCache.UpdateVenue(id, venue)
		if err != nil {
			errMsg := fmt.Sprintf("failed to update venue: %v", err)
			return nil, http.StatusInternalServerError, errors.New(errMsg)
		}
		err = s.SyncService.SyncVenueUpdate(id)
		if err != nil {
			errMsg := fmt.Sprintf("failed to sync change for venue: %v", err)
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
		s.SyncService.SyncVenueDelete(id)
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
		var artist domain.Artist
		if err := json.NewDecoder(r.Body).Decode(&artist); err != nil {
			return nil, http.StatusBadRequest, errors.New("invalid body")
		}
		savedArtist, err := s.ArtistCache.AddArtist(artist)
		if err != nil {
			errMsg := fmt.Sprintf("failed to save artist: %v", err)
			return nil, http.StatusInternalServerError, errors.New(errMsg)
		}
		err = s.SyncService.SyncArtistAdd(savedArtist.ID.Primary)
		if err != nil {
			errMsg := fmt.Sprintf("failed to sync change for artist: %v", err)
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
		var artist domain.Artist
		if err := json.NewDecoder(r.Body).Decode(&artist); err != nil {
			return nil, http.StatusBadRequest, errors.New("invalid body")
		}
		err := s.ArtistCache.UpdateArtist(id, artist)
		if err != nil {
			errMsg := fmt.Sprintf("failed to update artist: %v", err)
			return nil, http.StatusInternalServerError, errors.New(errMsg)
		}
		err = s.SyncService.SyncArtistUpdate(id)
		if err != nil {
			errMsg := fmt.Sprintf("failed to sync change for artist: %v", err)
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
		s.SyncService.SyncArtistDelete(id)
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
			if event.ID.Primary == id {
				return []domain.Event{event}, 0, nil
			}
		}
		errMsg := fmt.Sprintf("event with ID %s not found", id)
		return nil, http.StatusNotFound, errors.New(errMsg)
	case http.MethodPost:
		var event domain.Event
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			return nil, http.StatusBadRequest, errors.New("invalid body")
		}
		savedEvent, err := s.SavedEventCache.AddSavedEvent(event)
		if err != nil {
			errMsg := fmt.Sprintf("failed to save event: %v", err)
			return nil, http.StatusInternalServerError, errors.New(errMsg)
		}
		err = s.SyncService.SyncEventAdd(savedEvent.ID.Primary)
		if err != nil {
			errMsg := fmt.Sprintf("failed to sync change for event: %v", err)
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
		s.SyncService.SyncEventDelete(id)
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

func (s *Server) reloadGenres(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodGet {
		return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
	}

	queryParams := r.URL.Query()
	artists := queryParams["artists"]

	updateCount, err := s.ArtistInfoLoader.ReloadGenres(r.Context(), artists)
	if err != nil {
		errMsg := fmt.Sprintf("some artist genres failed to refresh, %s", err)
		return nil, http.StatusInternalServerError, errors.New(errMsg)
	}
	return fmt.Sprintf("Successfully refreshed genres for %d artists", updateCount), 0, nil
}

var thresholdOpts = map[string]ranker.RecLevel{
	"low":    ranker.LowMinRec,
	"medium": ranker.MediumMinRec,
	"high":   ranker.HighMinRec,
}

func (s *Server) getRecommendations(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodGet {
		return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
	}

	thresholdParam := r.URL.Query().Get("threshold")
	thresholdRfd := strings.ToLower(thresholdParam)
	threshold, exists := thresholdOpts[thresholdRfd]
	if !exists {
		errMsg := fmt.Sprintf("Invalid threshold: %s. Expected {low, medium, high}", thresholdParam)
		return nil, http.StatusBadRequest, errors.New(errMsg)
	}

	log.Info("Received GET recommendations request")
	recs := s.UpcomingEventsCache.GetRecommendedEvents(threshold)
	return recs, 0, nil
}

func (s *Server) getUpcomingEvents(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodGet {
		return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
	}

	events := s.UpcomingEventsCache.GetUpcomingEvents()
	return events, 0, nil
}

func (s *Server) refreshUpcomingEvents(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
	}

	err := s.UpcomingEventsCache.RefreshUpcomingEvents()
	if err != nil {
		log.Errorf("Failed to refresh upcoming events %v", err)
		return nil, http.StatusInternalServerError, errors.New("failed to refresh upcoming event cache")
	}
	return nil, 0, nil
}

func (s *Server) refreshRanks(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
	}

	log.Info("Manual ranks refresh triggered via API")
	go s.RanksCache.DoRefresh()
	return map[string]string{"status": "refresh started"}, http.StatusOK, nil
}

func (s *Server) handleGenres(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodGet {
		return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
	}

	genres := s.ArtistCache.GetUniqueGenres()
	return genres, 0, nil
}
