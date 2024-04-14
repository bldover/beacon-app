package spotify

import (
	"concert-manager/log"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	refreshUrl = "https://accounts.spotify.com/api/token"
	refreshParamsFmt = "?grant_type=refresh_token&refresh_token=%s"
	refreshKey = "SPOTIFY_REFRESH_TOKEN"
	authKey = "SPOTIFY_AUTH_TOKEN"
)

type authentication struct {
	accessToken string
    refreshToken string
	authToken string
}

func newAuthentication() *authentication {
	token := os.Getenv(refreshKey)
	if token == "" {
		log.Fatalf("%s env var must be set", refreshKey)
	}
	auth := os.Getenv(authKey)
	if auth == "" {
		log.Fatalf("%s env var must be set", authKey)
	}

	return &authentication{"", token, "Basic " + auth}
}

type refreshResponse struct {
    Token string `json:"access_token"`
	TokenType string `json:"token_type"`
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
	req.Header.Add("Authorization", a.authToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("received non-200 response from refresh call: %v, %s", resp.StatusCode, resp.Status)
		return errors.New(msg)
	}

	var token refreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return err
	}
	a.accessToken = token.TokenType + " " + token.Token
	return nil
}

// used for separately authenticating in a browser to get the initial refresh token
type SpotifyAuthHandler struct {}

func (h *SpotifyAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("Received Spotify authorization callback")
	out := "Request URI: " + r.RequestURI
	w.Write([]byte(out))
}
