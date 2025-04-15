package server

import (
	"concert-manager/finder"
	"concert-manager/log"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var thresholdOpts = map[string]finder.RecLevel{
	"low": finder.LowMinRec,
	"medium": finder.MediumMinRec,
	"high": finder.HighMinRec,
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
	recs := s.RecommendationCache.GetRecommendedEvents(threshold)
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
