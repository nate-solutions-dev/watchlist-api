package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL     string
	JWTSecret       string
	Port            string
	TMDBAccessToken string
	TMDBUsername    string
	TMDBPassword    string
	Environment     string
}

// Load reads environment variables and builds application config.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		Port:            strings.TrimSpace(os.Getenv("PORT")),
		TMDBAccessToken: os.Getenv("TMDB_ACCESS_TOKEN"),
		TMDBUsername:    os.Getenv("TMDB_USERNAME"),
		TMDBPassword:    os.Getenv("TMDB_PASSWORD"),
		Environment:     os.Getenv("APP_ENV"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.TMDBAccessToken == "" {
		return nil, fmt.Errorf("TMDB_ACCESS_TOKEN is required")
	}
	if cfg.Environment == "" {
		cfg.Environment = "development"
	}

	return cfg, nil
}
