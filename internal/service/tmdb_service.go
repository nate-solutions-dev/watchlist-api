package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Friel909/watchlist-api/config"
	"github.com/Friel909/watchlist-api/internal/logger"
)

type TMDBService interface {
	FetchMetadata(ctx context.Context, tmdbID int, mediaType string) (title string, posterURL string, genres []string, err error)
	GetSessionToken(ctx context.Context) (string, error)
}

type tmdbService struct {
	cfg *config.Config
}

// NewTMDBService creates a TMDB service backed by application config.
func NewTMDBService(cfg *config.Config) TMDBService {
	return &tmdbService{cfg: cfg}
}

type tmdbGenre struct {
	Name string `json:"name"`
}

type tmdbDetailResponse struct {
	Title      string      `json:"title"`
	Name       string      `json:"name"`
	PosterPath string      `json:"poster_path"`
	Genres     []tmdbGenre `json:"genres"`
}

type tmdbRequestTokenResponse struct {
	RequestToken string `json:"request_token"`
	Success      bool   `json:"success"`
}

type tmdbSessionResponse struct {
	SessionID string `json:"session_id"`
	Success   bool   `json:"success"`
}

type tmdbErrorResponse struct {
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
	Success       bool   `json:"success"`
}

// FetchMetadata gets title, poster url, and genres from TMDB detail endpoint.
func (s *tmdbService) FetchMetadata(ctx context.Context, tmdbID int, mediaType string) (string, string, []string, error) {
	if s.cfg.TMDBAccessToken == "" {
		return "", "", nil, fmt.Errorf("TMDB_ACCESS_TOKEN is required")
	}

	url := fmt.Sprintf("https://api.themoviedb.org/3/%s/%d", mediaType, tmdbID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", "", nil, fmt.Errorf("build tmdb request: %w", err)
	}

	var payload tmdbDetailResponse
	if err := s.doTMDBRequest(req, &payload); err != nil {
		return "", "", nil, err
	}

	title := payload.Title
	if title == "" {
		title = payload.Name
	}

	posterURL := ""
	if payload.PosterPath != "" {
		posterURL = "https://image.tmdb.org/t/p/w500" + payload.PosterPath
	}

	genres := make([]string, 0, len(payload.Genres))
	for _, genre := range payload.Genres {
		genres = append(genres, genre.Name)
	}

	return title, posterURL, genres, nil
}

// GetSessionToken creates a TMDB session token using configured TMDB credentials.
func (s *tmdbService) GetSessionToken(ctx context.Context) (string, error) {
	if s.cfg.TMDBAccessToken == "" {
		return "", fmt.Errorf("TMDB_ACCESS_TOKEN is required")
	}
	if s.cfg.TMDBUsername == "" || s.cfg.TMDBPassword == "" {
		return "", fmt.Errorf("TMDB_USERNAME and TMDB_PASSWORD are required")
	}

	requestToken, err := s.createRequestToken(ctx)
	if err != nil {
		return "", err
	}

	approvedToken, err := s.validateWithLogin(ctx, requestToken)
	if err != nil {
		return "", err
	}

	return s.createSession(ctx, approvedToken)
}

// createRequestToken requests a temporary TMDB request token.
func (s *tmdbService) createRequestToken(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.themoviedb.org/3/authentication/token/new", nil)
	if err != nil {
		return "", fmt.Errorf("build tmdb request token request: %w", err)
	}

	var respPayload tmdbRequestTokenResponse
	if err := s.doTMDBRequest(req, &respPayload); err != nil {
		return "", err
	}
	if !respPayload.Success || respPayload.RequestToken == "" {
		return "", fmt.Errorf("tmdb did not return a valid request token")
	}

	return respPayload.RequestToken, nil
}

// validateWithLogin validates TMDB username/password against request token.
func (s *tmdbService) validateWithLogin(ctx context.Context, requestToken string) (string, error) {
	body, err := json.Marshal(map[string]string{
		"username":      s.cfg.TMDBUsername,
		"password":      s.cfg.TMDBPassword,
		"request_token": requestToken,
	})
	if err != nil {
		return "", fmt.Errorf("marshal tmdb login payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.themoviedb.org/3/authentication/token/validate_with_login", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("build tmdb login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	var respPayload tmdbRequestTokenResponse
	if err := s.doTMDBRequest(req, &respPayload); err != nil {
		return "", err
	}
	if !respPayload.Success || respPayload.RequestToken == "" {
		return "", fmt.Errorf("tmdb did not validate login credentials")
	}

	return respPayload.RequestToken, nil
}

// createSession exchanges a validated request token for TMDB session id.
func (s *tmdbService) createSession(ctx context.Context, requestToken string) (string, error) {
	body, err := json.Marshal(map[string]string{
		"request_token": requestToken,
	})
	if err != nil {
		return "", fmt.Errorf("marshal tmdb session payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.themoviedb.org/3/authentication/session/new", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("build tmdb session request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	var respPayload tmdbSessionResponse
	if err := s.doTMDBRequest(req, &respPayload); err != nil {
		return "", err
	}
	if !respPayload.Success || respPayload.SessionID == "" {
		return "", fmt.Errorf("tmdb did not return a valid session id")
	}

	return respPayload.SessionID, nil
}

// doTMDBRequest performs a TMDB HTTP call and decodes JSON response payload.
func (s *tmdbService) doTMDBRequest(req *http.Request, target any) error {
	req.Header.Set("Authorization", "Bearer "+s.cfg.TMDBAccessToken)
	req.Header.Set("Accept", "application/json")
	start := time.Now()

	logger.Info(req.Context(), "TMDBClient.DoRequest", "sending request", "method", req.Method, "url", req.URL.String())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(req.Context(), "TMDBClient.DoRequest", "request failed", "url", req.URL.String(), "duration", time.Since(start).String(), "error", err.Error())
		return fmt.Errorf("call tmdb: %w", err)
	}
	defer resp.Body.Close()
	logger.Info(req.Context(), "TMDBClient.DoRequest", "request completed", "url", req.URL.String(), "status", resp.StatusCode, "duration", time.Since(start).String())

	if resp.StatusCode >= 400 {
		var tmdbErr tmdbErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&tmdbErr); err == nil && tmdbErr.StatusMessage != "" {
			logger.Warn(req.Context(), "TMDBClient.DoRequest", "tmdb returned error", "status", tmdbErr.StatusCode)
			return fmt.Errorf("tmdb error %d: %s", tmdbErr.StatusCode, tmdbErr.StatusMessage)
		}
		logger.Warn(req.Context(), "TMDBClient.DoRequest", "tmdb returned error", "status", resp.StatusCode)
		return fmt.Errorf("tmdb returned status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		logger.Error(req.Context(), "TMDBClient.DoRequest", "decode failed", "error", err.Error())
		return fmt.Errorf("decode tmdb response: %w", err)
	}

	return nil
}
