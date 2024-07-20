package server

import (
	"concert-manager/cache"
	"concert-manager/data"
	"concert-manager/log"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Server struct {
    Loader loader
	SavedEventCache savedEventCache
	ArtistCache artistCache
	VenueCache venueCache
	UpcomingEventsCache upcomingEventsCache
	RecommendationCache recommendationCache
}

type loader interface {
    Upload(context.Context, io.ReadCloser) (int, error)
}

type savedEventCache interface {
    GetSavedEvents() []data.Event
	GetPassedSavedEvents() []data.Event
	AddSavedEvent(data.Event) (*data.Event, error)
	DeleteSavedEvent(string) error
}

type artistCache interface {
    GetArtists() []data.Artist
	AddArtist(data.Artist) (*data.Artist, error)
	DeleteArtist(string) error
}

type venueCache interface {
    GetVenues() []data.Venue
	AddVenue(data.Venue) (*data.Venue, error)
	DeleteVenue(string) error
}

type upcomingEventsCache interface {
    GetUpcomingEvents() []data.EventDetails
	ChangeLocation(string, string)
	GetLocation() cache.Location
}

type recommendationCache interface {
    GetRecommendedEvents(cache.Threshold) []data.EventRank
}

const port = ":3001"

func (s *Server) StartServer() {
	http.HandleFunc("/v1/upload", s.handleRequest(s.handleUpload))
	http.HandleFunc("/v1/events/upcoming", s.handleRequest(s.getUpcomingEvents))
	http.HandleFunc("/v1/events/recommended", s.handleRequest(s.getRecommendations))
	http.HandleFunc("/v1/events/saved", s.handleRequest(s.handleSavedEvents))
	http.HandleFunc("/v1/venues", s.handleRequest(s.handleVenues))
	http.HandleFunc("/v1/artists", s.handleRequest(s.handleArtists))
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
