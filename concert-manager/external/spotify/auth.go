package spotify

import (
	"concert-manager/file"
	"concert-manager/log"
	"crypto/rand"
	"encoding/base64"
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
	authorizeUrl             = "https://accounts.spotify.com/authorize"
	tokenUrl                 = "https://accounts.spotify.com/api/token"
	clientIdKey              = "CM_SPOTIFY_CLIENT_ID"
	clientSecretKey          = "CM_SPOTIFY_CLIENT_SECRET"
	callbackUrlKey           = "CM_SPOTIFY_CALLBACK_URL"
	refreshTokenFile         = "spotify_refresh_token"
	refreshTokenDurationDays = 180 * 24 * time.Hour
	authScopes               = "playlist-read-private playlist-read-collaborative user-top-read user-library-read"
	reauthStateTTL           = 10 * time.Minute
	reauthStateByteLength    = 32
)

type authentication struct {
	RefreshTokenExpireTs time.Time
	clientId             string
	authKey              string
	callbackUrl          string
	reauthStateStore     []string
	reauthStateMutex     sync.Mutex
	tokenMutex           sync.Mutex
	refreshToken         string
	accessToken          string
	accessTokenInvalid   []string
}

type storedRefreshToken struct {
	Token    string    `json:"token"`
	ExpireTs time.Time `json:"expireTs"`
}

func NewAuthentication() *authentication {
	clientId := os.Getenv(clientIdKey)
	if clientId == "" {
		log.Fatalf("%s env var must be set", clientIdKey)
	}
	clientSecret := os.Getenv(clientSecretKey)
	if clientSecret == "" {
		log.Fatalf("%s env var must be set", clientSecretKey)
	}
	authToken := buildAuthToken(clientId, clientSecret)

	callbackUrl := os.Getenv(callbackUrlKey)
	if callbackUrl == "" {
		log.Fatalf("%s env var must be set", callbackUrlKey)
	}

	stored, err := loadRefreshToken()
	if err != nil {
		log.Alertf("Failed to load Spotify refresh token. Reauthentication is required, %v", err)
		stored = storedRefreshToken{}
	} else {
		log.Infof("Loaded Spotify refresh token")
	}

	auth := &authentication{
		clientId:             clientId,
		authKey:              authToken,
		callbackUrl:          callbackUrl,
		reauthStateStore:     []string{},
		reauthStateMutex:     sync.Mutex{},
		tokenMutex:           sync.Mutex{},
		refreshToken:         stored.Token,
		RefreshTokenExpireTs: stored.ExpireTs,
		accessToken:          "",
		accessTokenInvalid:   []string{},
	}

	if auth.refreshToken != "" {
		// checks refresh token validity
		auth.refreshAccessToken()
	}

	return auth
}

func (a *authentication) GetAuthStatus() (bool, time.Time) {
	a.tokenMutex.Lock()
	defer a.tokenMutex.Unlock()
	authenticated := a.refreshToken != "" && time.Now().Before(a.RefreshTokenExpireTs)
	return authenticated, a.RefreshTokenExpireTs
}

func (a *authentication) StartReauth() (string, error) {
	state, err := a.issueReauthState()
	if err != nil {
		return "", fmt.Errorf("failed to generate OAuth state: %w", err)
	}

	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", a.clientId)
	params.Set("redirect_uri", a.callbackUrl)
	params.Set("scope", authScopes)
	params.Set("state", state)
	return authorizeUrl + "?" + params.Encode(), nil
}

func (a *authentication) CompleteReauth(code, state string) error {
	if strings.TrimSpace(code) == "" {
		return errors.New("missing authorization code")
	}
	if !a.consumeState(state) {
		return errors.New("invalid or expired state")
	}
	return a.exchangeCode(code)
}

func buildAuthToken(clientId, clientSecret string) string {
	authToken := []byte(clientId + ":" + clientSecret)
	return "Basic " + base64.StdEncoding.EncodeToString(authToken)
}

func loadRefreshToken() (storedRefreshToken, error) {
	path, err := file.GetCacheFilePath(refreshTokenFile)
	if err != nil {
		return storedRefreshToken{}, err
	}

	if !file.FileExists(path) {
		return storedRefreshToken{}, errors.New("refresh token file does not exist")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return storedRefreshToken{}, err
	}

	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" {
		return storedRefreshToken{}, errors.New("refresh token file is empty")
	}

	var record storedRefreshToken
	if err := json.Unmarshal([]byte(trimmed), &record); err != nil {
		return storedRefreshToken{}, err
	}

	return record, nil
}

func persistRefreshToken(record storedRefreshToken) error {
	path, err := file.GetCacheFilePath(refreshTokenFile)
	if err != nil {
		return err
	}
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}

func (a *authentication) getAccessToken() (string, error) {
	a.tokenMutex.Lock()
	defer a.tokenMutex.Unlock()
	if a.accessToken == "" {
		if err := a.refreshAccessToken(); err != nil {
			errMsg := fmt.Sprintf("failed to refresh Spotify access token: %v", err)
			return "", errors.New(errMsg)
		}
	}
	return a.accessToken, nil
}

type refreshResponse struct {
	Token     string `json:"access_token"`
	TokenType string `json:"token_type"`
}

type authErrorResponse struct {
	Error     string `json:"error"`
	ErrorDesc string `json:"error_description"`
}

func (a *authentication) refreshAccessToken() error {
	log.Debug("Refreshing Spotify access token")
	if a.refreshToken == "" {
		return errors.New("invalid refresh token")
	}

	params := url.Values{}
	params.Add("grant_type", "refresh_token")
	params.Add("refresh_token", a.refreshToken)
	req, err := http.NewRequest(http.MethodPost, tokenUrl, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", a.authKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("received non-200 response from access token refresh call: %v, %s", resp.StatusCode, resp.Status)
		var errResp authErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil {
			errMsg += fmt.Sprintf(", %+v", errResp)
			if errResp.Error == "invalid_grant" {
				a.refreshToken = ""
			}
		}
		return errors.New(errMsg)
	}

	var apiToken refreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiToken); err != nil {
		return err
	}
	a.accessToken = apiToken.TokenType + " " + apiToken.Token
	return nil
}

func (a *authentication) markAccessTokenExpired(token string) {
	a.tokenMutex.Lock()
	defer a.tokenMutex.Unlock()
	if !slices.Contains(a.accessTokenInvalid, token) {
		// this accessTokenInvalid list doesn't seem needed? consider removing it
		a.accessTokenInvalid = append(a.accessTokenInvalid, token)
		a.accessToken = ""
		// give more than enough time for any parallel calls to finish with the expired token
		cleanupDelay, _ := time.ParseDuration("1m")
		time.AfterFunc(cleanupDelay, func() {
			a.tokenMutex.Lock()
			defer a.tokenMutex.Unlock()
			slices.DeleteFunc(a.accessTokenInvalid, func(t string) bool {
				return t == token
			})
		})
	}
}

func (a *authentication) issueReauthState() (string, error) {
	buf := make([]byte, reauthStateByteLength)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	state := base64.RawURLEncoding.EncodeToString(buf)

	a.reauthStateMutex.Lock()
	defer a.reauthStateMutex.Unlock()
	a.reauthStateStore = append(a.reauthStateStore, state)
	time.AfterFunc(reauthStateTTL, func() {
		a.reauthStateMutex.Lock()
		defer a.reauthStateMutex.Unlock()
		slices.DeleteFunc(a.reauthStateStore, func(s string) bool {
			return s == state
		})
	})
	return state, nil
}

func (a *authentication) consumeState(state string) bool {
	a.reauthStateMutex.Lock()
	defer a.reauthStateMutex.Unlock()
	if !slices.Contains(a.reauthStateStore, state) {
		return false
	}

	slices.DeleteFunc(a.reauthStateStore, func(s string) bool {
		return s == state
	})
	return true
}

type codeExchangeResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

func (a *authentication) exchangeCode(code string) error {
	params := url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("code", code)
	params.Add("redirect_uri", a.callbackUrl)
	req, err := http.NewRequest(http.MethodPost, tokenUrl, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", a.authKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("received non-200 response from code exchange: %v, %s", resp.StatusCode, resp.Status)
		var errResp authErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil {
			errMsg += fmt.Sprintf(", %+v", errResp)
		}
		return errors.New(errMsg)
	}

	var body codeExchangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return err
	}
	if body.RefreshToken == "" {
		return errors.New("code exchange succeeded but response contained no refresh_token")
	}

	a.tokenMutex.Lock()
	defer a.tokenMutex.Unlock()
	expireTs := time.Now().Add(refreshTokenDurationDays)
	record := storedRefreshToken{Token: body.RefreshToken, ExpireTs: expireTs}
	if err := persistRefreshToken(record); err != nil {
		log.Alertf("Failed to persist new Spotify refresh token: %w", err)
	}
	a.refreshToken = body.RefreshToken
	a.RefreshTokenExpireTs = expireTs
	a.accessToken = body.TokenType + " " + body.AccessToken
	log.Info("Successfully rotated Spotify refresh token")
	return nil
}
