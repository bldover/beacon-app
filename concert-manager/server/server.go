package server

import (
	"concert-manager/data"
	"concert-manager/finder"
	"concert-manager/log"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Server struct {
    EventLoader eventLoader
	ArtistInfoLoader artistInfoLoader
	SavedEventCache savedEventCache
	ArtistCache artistCache
	VenueCache venueCache
	UpcomingEventsCache upcomingEventsCache
	RecommendationCache recommendationCache
}

type eventLoader interface {
    Upload(context.Context, io.ReadCloser) (int, error)
}

type artistInfoLoader interface {
    PopulateLastFmData(context.Context, []string) (int, error)
	PopulateUserGenres(context.Context) (int, error)
}

type savedEventCache interface {
    GetSavedEvents() []data.Event
	GetPassedSavedEvents() []data.Event
	AddSavedEvent(data.Event) (*data.Event, error)
	DeleteSavedEvent(string) error
	RefreshSavedEvents() error
}

type artistCache interface {
    GetArtists() []data.Artist
	AddArtist(data.Artist) (*data.Artist, error)
	UpdateArtist(string, data.Artist) error
	DeleteArtist(string) error
	RefreshArtists() error
}

type venueCache interface {
    GetVenues() []data.Venue
	AddVenue(data.Venue) (*data.Venue, error)
	UpdateVenue(string, data.Venue) error
	DeleteVenue(string) error
	RefreshVenues() error
}

type upcomingEventsCache interface {
    GetUpcomingEvents() []data.EventDetails
	ChangeLocation(string, string)
	GetLocation() finder.Location
	RefreshUpcomingEvents() error
}

type recommendationCache interface {
    GetRecommendedEvents(finder.RecLevel) []data.EventDetails
}

const port = ":3001"

// TODO: Use (or write?) a better HTTP library to clean up the server routing and handlers
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
	http.HandleFunc("/v1/artists/loadlastfm", s.handleRequest(s.loadLastFmGenres))
	http.HandleFunc("/v1/genres/loaduser", s.handleRequest(s.loadUserGenres))
//	http.Handle("/spotify/callback", &spotify.SpotifyAuthHandler{})

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
