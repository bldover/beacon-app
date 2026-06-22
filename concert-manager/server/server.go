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
	"slices"
	"strings"
	"time"
)

type Server struct {
	EventLoader         eventLoader
	ArtistInfoLoader    artistInfoLoader
	SavedEventCache     savedEventStore
	ArtistCache         artistStore
	VenueCache          venueStore
	AlbumCache          albumStore
	UpcomingEventsCache upcomingEventsStore
	RanksCache          ranksRefresher
	SyncService         dataSyncService
	ImageUploader       imageUploader
	SpotifyAuthHandler  spotifyOAuthHandler
	ApiKey              string
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
	GetUniqueGenres() domain.GenreResponse
}

type venueStore interface {
	GetVenues() []domain.Venue
	AddVenue(domain.Venue) (*domain.Venue, error)
	UpdateVenue(string, domain.Venue) error
	DeleteVenue(string) error
	RefreshVenues() error
}

type albumStore interface {
	GetAlbums() []domain.Album
	AddAlbum(domain.Album) (*domain.Album, error)
	UpdateAlbum(string, domain.Album) error
	DeleteAlbum(string) error
	RefreshAlbums() error
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

type ranksRefresher interface {
	DoRefresh()
}

type imageUploader interface {
	UploadImage(context.Context, io.Reader, string) (string, error)
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
	http.HandleFunc("/v1/albums", s.handleRequest(s.handleAlbums))
	http.HandleFunc("/v1/albums/", s.handleRequest(s.handleAlbums))
	http.HandleFunc("/v1/albums/refresh", s.handleRequest(s.refreshAlbums))
	http.HandleFunc("/v1/albums/images", s.handleRequest(s.handleAlbumImages))
	http.HandleFunc("/v1/artists", s.handleRequest(s.handleArtists))
	http.HandleFunc("/v1/artists/", s.handleRequest(s.handleArtists))
	http.HandleFunc("/v1/artists/refresh", s.handleRequest(s.refreshArtists))
	http.HandleFunc("/v1/ranks/refresh", s.handleRequest(s.refreshRanks))
	http.HandleFunc("/v1/genres", s.handleRequest(s.handleGenres))
	http.HandleFunc("/v1/genres/refresh", s.handleRequest(s.reloadGenres))
	http.HandleFunc("/v1/analytics/summary", s.handleRequest(s.getAnalyticsSummary))
	http.HandleFunc("/v1/analytics/years", s.handleRequest(s.handleAnalyticsYears))
	http.HandleFunc("/v1/analytics/years/", s.handleRequest(s.handleAnalyticsYears))
	http.HandleFunc("/v1/analytics/months", s.handleRequest(s.handleAnalyticsMonths))
	http.HandleFunc("/v1/analytics/months/", s.handleRequest(s.handleAnalyticsMonths))
	http.HandleFunc("/v1/analytics/artists", s.handleRequest(s.handleAnalyticsArtists))
	http.HandleFunc("/v1/analytics/artists/", s.handleRequest(s.handleAnalyticsArtists))
	http.HandleFunc("/v1/analytics/venues", s.handleRequest(s.handleAnalyticsVenues))
	http.HandleFunc("/v1/analytics/venues/", s.handleRequest(s.handleAnalyticsVenues))
	http.HandleFunc("/v1/analytics/genres", s.handleRequest(s.handleAnalyticsGenres))
	http.HandleFunc("/v1/analytics/genres/", s.handleRequest(s.handleAnalyticsGenres))
	http.HandleFunc("/v1/spotify/auth/start", s.handleRequest(s.startSpotifyAuth))
	http.HandleFunc("/v1/spotify/auth/status", s.handleRequest(s.getSpotifyAuthStatus))
	// doesn't use handleRequest for custom deep-link response
	http.HandleFunc("/v1/spotify/auth/callback", s.handleSpotifyAuthCallback)

	log.Info("Starting server on port", port)
	log.Fatal(http.ListenAndServe(port, s.authMiddleware(http.DefaultServeMux)))
}

// unauthenticated for OAuth callback
var publicPaths = map[string]bool{
	"/v1/spotify/auth/callback": true,
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if publicPaths[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}
		const prefix = "Bearer "
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, prefix) {
			log.Infof("Unauthorized request to %s from %s: missing bearer token", r.URL.Path, r.RemoteAddr)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		got := strings.TrimPrefix(header, prefix)
		if slices.Compare([]byte(got), []byte(s.ApiKey)) != 0 {
			log.Infof("Unauthorized request to %s from %s: invalid bearer token", r.URL.Path, r.RemoteAddr)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
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
