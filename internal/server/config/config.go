package config

import (
	"fmt"
	"os"
	"time"
)

// Config via .env files
type Config struct {
	Env          string        // APP_ENV
	DatabaseDSN  string        // DATABASE_DSN
	AppSecret    string        // APP_SECRET
	Port         string        // PORT
	ReadTimeout  time.Duration // READ_TIMEOUT
	WriteTimeout time.Duration // WRITE_TIMEOUT
	Debug        bool          // DEBUG
}

const (
	defaultPort         = ":50051"
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 10 * time.Second
)

func NewConfig() (*Config, error) {
	cfg := &Config{
		Env:          getEnv("APP_ENV", "dev"),
		DatabaseDSN:  os.Getenv("DATABASE_DSN"),
		AppSecret:    os.Getenv("APP_SECRET"),
		Port:         getEnv("PORT", defaultPort),
		ReadTimeout:  parseDuration("READ_TIMEOUT", defaultReadTimeout),
		WriteTimeout: parseDuration("WRITE_TIMEOUT", defaultWriteTimeout),
		Debug:        getBoolEnv("DEBUG"),
	}

	if cfg.DatabaseDSN == "" {
		return nil, fmt.Errorf("DATABASE_DSN is required")
	}

	if cfg.AppSecret == "" {
		return nil, fmt.Errorf("APP_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getBoolEnv(key string) bool {
	return os.Getenv(key) == "true" || os.Getenv(key) == "1"
}

func parseDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return duration
}
