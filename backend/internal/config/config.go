package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var (
	ErrMissingDatabaseURL = errors.New("DATABASE_URL is required")
	ErrMissingJWTSecret   = errors.New("JWT_SECRET is required")
)

type Config struct {
	AppEnv      string
	Port        string
	DatabaseURL string
	JWTSecret   string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func (c Config) Validate() error {
	if c.DatabaseURL == "" {
		return ErrMissingDatabaseURL
	}

	if c.JWTSecret == "" {
		return ErrMissingJWTSecret
	}

	return nil
}
