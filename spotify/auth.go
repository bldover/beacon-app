package spotify

import (
	"concert-manager/log"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"sync"
	"time"
)

const (
	refreshUrl = "https://accounts.spotify.com/api/token"
	refreshParamsFmt = "?grant_type=refresh_token&refresh_token=%s"
	refreshKey = "CM_SPOTIFY_REFRESH_TOKEN"
	authKey = "CM_SPOTIFY_AUTH_TOKEN"
)

type authentication struct {
	authMutex sync.Mutex
	invalid []string
    refreshToken string
	refreshAuthHeader string
	apiAuthToken string
}

func newAuthentication() *authentication {
	refreshToken := os.Getenv(refreshKey)
	if refreshToken == "" {
		log.Fatalf("%s env var must be set", refreshKey)
	}
	auth := os.Getenv(authKey)
	if auth == "" {
		log.Fatalf("%s env var must be set", authKey)
	}

	return &authentication{
		authMutex: sync.Mutex{},
		invalid: []string{},
		refreshToken: refreshToken,
		refreshAuthHeader: "Basic " + auth,
		apiAuthToken: "",
	}
}

type refreshResponse struct {
    Token string `json:"access_token"`
	TokenType string `json:"token_type"`
}

func (a *authentication) getAuthToken() (string, error) {
	a.authMutex.Lock()
	defer a.authMutex.Unlock()
	if a.apiAuthToken == "" {
		if err := a.refresh(); err != nil {
			errMsg := fmt.Sprintf("failed to refresh Spotify token: %v", err)
			return "", errors.New(errMsg)
		}
	}
	return a.apiAuthToken, nil
}

func (a *authentication) markAuthExpired(token string) {
	a.authMutex.Lock()
	defer a.authMutex.Unlock()
	if !slices.Contains(a.invalid, token) {
		a.invalid = append(a.invalid, token)
		a.apiAuthToken = ""
		// give more than enough time for any parallel calls to finish with the expired token
		cleanupDelay, _ := time.ParseDuration("1m")
		time.AfterFunc(cleanupDelay, func() {
			a.authMutex.Lock()
			defer a.authMutex.Unlock()
			slices.DeleteFunc(a.invalid, func(t string) bool {
				return t == token
			})
		})
	}
}

func (a *authentication) refresh() error {
	log.Debug("Refreshing Spotify access token")
	params := url.Values{}
	params.Add("grant_type", "refresh_token")
	params.Add("refresh_token", a.refreshToken)
	req, err := http.NewRequest(http.MethodPost, refreshUrl, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", a.refreshAuthHeader)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("received non-200 response from refresh call: %v, %s", resp.StatusCode, resp.Status)
		return errors.New(msg)
	}

	var apiToken refreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiToken); err != nil {
		return err
	}
	a.apiAuthToken = apiToken.TokenType + " " + apiToken.Token
	return nil
}

// used for separately authenticating in a browser to get the initial refresh token
type SpotifyAuthHandler struct {}

func (h *SpotifyAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("Received Spotify authorization callback")
	out := "Request URI: " + r.RequestURI
	w.Write([]byte(out))
}
