// Package config reads application settings from environment variables.
package config

import (
	"errors"
	"fmt"
)

// Config holds every setting the server needs. It is built once at startup
// and passed down explicitly; nothing reads the environment after that.
type Config struct {
	DatabaseURL string
	HTTPAddr    string
	LogFormat   string // "json" or "text"
}

// Load builds a Config using getenv to look up variables. Production code
// passes os.Getenv; tests pass a fake backed by a map.
//
// All problems are collected and returned together, so a misconfigured
// deploy reports every missing variable at once instead of one per restart.
func Load(getenv func(string) string) (Config, error) {
	cfg := Config{
		DatabaseURL: getenv("DATABASE_URL"),
		HTTPAddr:    withDefault(getenv("HTTP_ADDR"), ":8080"),
		LogFormat:   withDefault(getenv("LOG_FORMAT"), "json"),
	}

	var errs []error
	if cfg.DatabaseURL == "" {
		errs = append(errs, errors.New("DATABASE_URL is required"))
	}
	if cfg.LogFormat != "json" && cfg.LogFormat != "text" {
		errs = append(errs, fmt.Errorf("LOG_FORMAT must be json or text, got %q", cfg.LogFormat))
	}

	if err := errors.Join(errs...); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func withDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
