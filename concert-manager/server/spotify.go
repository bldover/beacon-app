package server

import (
	"concert-manager/log"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const spotifyAuthDeepLink = "beacon://spotify-auth-complete"

type spotifyOAuthHandler interface {
	StartReauth() (string, error)
	CompleteReauth(code, state string) error
	GetAuthStatus() (bool, time.Time)
}

type spotifyReauthResponse struct {
	AuthUrl string `json:"authUrl"`
}

type spotifyAuthStatusResponse struct {
	Authenticated bool      `json:"authenticated"`
	ExpireTs      time.Time `json:"expireTs"`
}

func (s *Server) startSpotifyAuth(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodGet {
		return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
	}
	authUrl, err := s.SpotifyAuthHandler.StartReauth()
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to start Spotify reauth: %w", err)
	}
	return spotifyReauthResponse{AuthUrl: authUrl}, 0, nil
}

func (s *Server) getSpotifyAuthStatus(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodGet {
		return nil, http.StatusMethodNotAllowed, errors.New("unsupported method")
	}
	authenticated, expireTs := s.SpotifyAuthHandler.GetAuthStatus()
	return spotifyAuthStatusResponse{Authenticated: authenticated, ExpireTs: expireTs}, 0, nil
}

func (s *Server) handleSpotifyAuthCallback(w http.ResponseWriter, r *http.Request) {
	log.Infof("Received Spotify callback from %s", r.RemoteAddr)
	q := r.URL.Query()

	if errCode := q.Get("error"); errCode != "" {
		log.Errorf("Spotify reported error during reauth: %s", errCode)
		redirectToApp(w, r, "error", errCode)
		return
	}

	code := q.Get("code")
	state := q.Get("state")
	if err := s.SpotifyAuthHandler.CompleteReauth(code, state); err != nil {
		log.Errorf("Spotify reauth failed: %v", err)
		redirectToApp(w, r, "error", err.Error())
		return
	}

	redirectToApp(w, r, "ok", "")
}

func redirectToApp(w http.ResponseWriter, r *http.Request, status, reason string) {
	params := url.Values{}
	params.Set("status", status)
	if reason != "" {
		params.Set("reason", reason)
	}
	http.Redirect(w, r, spotifyAuthDeepLink+"?"+params.Encode(), http.StatusFound)
}
