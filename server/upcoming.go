package server

import (
	"concert-manager/cache"
	"concert-manager/log"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var thresholdOpts = map[string]cache.Threshold{
	"low": cache.LowThreshold,
	"medium": cache.MediumThreshold,
	"high": cache.HighThreshold,
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
