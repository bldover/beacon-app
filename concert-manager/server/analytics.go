package server

import (
	"concert-manager/analytics"
	"concert-manager/domain"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

func (s *Server) pastEvents() []domain.Event {
	return analytics.PastEvents(s.SavedEventCache.GetSavedEvents())
}

func (s *Server) getAnalyticsSummary(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodGet {
		return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
	}
	return analytics.BuildSummary(s.pastEvents()), 0, nil
}

func (s *Server) handleAnalyticsYears(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodGet {
		return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
	}
	events := s.pastEvents()
	key, err := analyticsPathKey(r.URL.Path, "years")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	if key == "" {
		return analytics.CountByYear(events), 0, nil
	}
	filtered := analytics.FilterEventsByYear(events, key)
	return analytics.EventsResponse{Count: len(filtered), Events: filtered}, 0, nil
}

func (s *Server) handleAnalyticsMonths(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodGet {
		return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
	}
	events := s.pastEvents()
	key, err := analyticsPathKey(r.URL.Path, "months")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	if key == "" {
		return analytics.CountByMonth(events), 0, nil
	}
	filtered := analytics.FilterEventsByMonth(events, key)
	return analytics.EventsResponse{Count: len(filtered), Events: filtered}, 0, nil
}

func (s *Server) handleAnalyticsArtists(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodGet {
		return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
	}
	events := s.pastEvents()
	key, err := analyticsPathKey(r.URL.Path, "artists")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	if key == "" {
		return analytics.CountByArtist(events), 0, nil
	}
	filtered := analytics.FilterEventsByArtist(events, key)
	return analytics.EventsResponse{Count: len(filtered), Events: filtered}, 0, nil
}

func (s *Server) handleAnalyticsVenues(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodGet {
		return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
	}
	events := s.pastEvents()
	key, err := analyticsPathKey(r.URL.Path, "venues")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	if key == "" {
		return analytics.CountByVenue(events), 0, nil
	}
	filtered := analytics.FilterEventsByVenue(events, key)
	return analytics.EventsResponse{Count: len(filtered), Events: filtered}, 0, nil
}

func (s *Server) handleAnalyticsGenres(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodGet {
		return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
	}
	events := s.pastEvents()
	key, err := analyticsPathKey(r.URL.Path, "genres")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	if key == "" {
		return analytics.CountByGenre(events), 0, nil
	}
	filtered := analytics.FilterEventsByGenre(events, key)
	return analytics.EventsResponse{Count: len(filtered), Events: filtered}, 0, nil
}

// analyticsPathKey returns the optional final segment of /v1/analytics/{dim}/{key}.
// Returns ("", nil) when no key segment is present (the listing route).
func analyticsPathKey(path, dim string) (string, error) {
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	// expected /v1/analytics/{dim} or /v1/analytics/{dim}/{key}
	for i, p := range parts {
		if p == dim && i == len(parts)-1 {
			return "", nil
		}
		if p == dim && i < len(parts)-1 {
			rawKey := strings.Join(parts[i+1:], "/")
			decoded, err := url.PathUnescape(rawKey)
			if err != nil {
				return "", errors.New("invalid path encoding")
			}
			if decoded == "" {
				return "", nil
			}
			return decoded, nil
		}
	}
	return "", errors.New("malformed analytics path")
}
