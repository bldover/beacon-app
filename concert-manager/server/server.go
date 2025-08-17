package server

import (
	"concert-manager/domain"
	"concert-manager/finder"
	"concert-manager/log"
	"concert-manager/ranker"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Server struct {
	EventLoader         eventLoader
	ArtistInfoLoader    artistInfoLoader
	SavedEventCache     savedEventStore
	ArtistCache         artistStore
	VenueCache          venueStore
	UpcomingEventsCache upcomingEventsStore
	SyncService         dataSyncService
}

type eventLoader interface {
	Upload(context.Context, io.ReadCloser) (int, error)
}

type artistInfoLoader interface {
	ReloadGenres(context.Context, []string) (int, error)
}

type savedEventStore interface {
	GetSavedEvents() []domain.Event
	GetPassedSavedEvents() []domain.Event
	AddSavedEvent(domain.Event) (*domain.Event, error)
	DeleteSavedEvent(string) error
	RefreshSavedEvents() error
}

type artistStore interface {
	GetArtists() []domain.Artist
	AddArtist(domain.Artist) (*domain.Artist, error)
	UpdateArtist(string, domain.Artist) error
	DeleteArtist(string) error
	RefreshArtists() error
}

type venueStore interface {
	GetVenues() []domain.Venue
	AddVenue(domain.Venue) (*domain.Venue, error)
	UpdateVenue(string, domain.Venue) error
	DeleteVenue(string) error
	RefreshVenues() error
}

type upcomingEventsStore interface {
	GetUpcomingEvents() []domain.EventDetails
	ChangeLocation(string, string)
	GetLocation() finder.Location
	RefreshUpcomingEvents() error
	GetRecommendedEvents(ranker.RecLevel) []domain.EventDetails
}

type dataSyncService interface {
	SyncArtistAdd(string) error
	SyncArtistUpdate(string) error
	SyncArtistDelete(string)
	SyncVenueAdd(string) error
	SyncVenueUpdate(string) error
	SyncVenueDelete(string)
	SyncEventAdd(string) error
	SyncEventUpdate(string) error
	SyncEventDelete(string)
}

const port = ":3001"

func (s *Server) StartServer() {
	http.HandleFunc("/v1/upload", s.handleRequest(s.handleUpload))
	http.HandleFunc("/v1/events/upcoming", s.handleRequest(s.getUpcomingEvents))
	http.HandleFunc("/v1/events/upcoming/refresh", s.handleRequest(s.refreshUpcomingEvents))
	http.HandleFunc("/v1/events/recommended", s.handleRequest(s.getRecommendations))
	http.HandleFunc("/v1/events/saved", s.handleRequest(s.handleSavedEvents))
	http.HandleFunc("/v1/events/saved/", s.handleRequest(s.handleSavedEvents))
	http.HandleFunc("/v1/events/saved/refresh", s.handleRequest(s.refreshSavedEvents))
	http.HandleFunc("/v1/venues", s.handleRequest(s.handleVenues))
	http.HandleFunc("/v1/venues/", s.handleRequest(s.handleVenues))
	http.HandleFunc("/v1/venues/refresh", s.handleRequest(s.refreshVenues))
	http.HandleFunc("/v1/artists", s.handleRequest(s.handleArtists))
	http.HandleFunc("/v1/artists/", s.handleRequest(s.handleArtists))
	http.HandleFunc("/v1/artists/refresh", s.handleRequest(s.refreshArtists))
	http.HandleFunc("/v1/genres/refresh", s.handleRequest(s.reloadGenres))
	http.Handle("/auth/callback", &authHandler{})

	log.Info("Starting server on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

type handlerFunc func(http.ResponseWriter, *http.Request) (any, int, error)

func (s *Server) handleRequest(f handlerFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := time.Now().Nanosecond()
		log.Infof("Received request (%s) %s, assigned ID: %d", r.Method, r.URL, id)
		startTs := time.Now()
		body, status, err := f(w, r)
		if err != nil {
			log.Errorf("Error processing request ID %d: %v", id, err)
			http.Error(w, err.Error(), status)
		}
		if body != nil {
			json.NewEncoder(w).Encode(body)
		}
		log.Infof("Finished processing request ID %d in %v ms", id, time.Since(startTs).Milliseconds())
	}
}

// callback URL for retrieving OAUTH access tokens
type authHandler struct{}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("Received authorization callback")
	out := "Request URI: " + r.RequestURI
	w.Write([]byte(out))
}
