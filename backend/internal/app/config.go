package app

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr           string
	DBPath         string
	FrontendOrigin string
	IsProduction   bool
	SessionTTL     time.Duration
}

func LoadConfig() Config {
	return Config{
		Addr:           getenv("API_ADDR", ":8080"),
		DBPath:         getenv("DB_PATH", "./finance.db"),
		FrontendOrigin: getenv("FRONTEND_ORIGIN", "http://localhost:5173"),
		IsProduction:   getenvBool("APP_PRODUCTION", false),
		SessionTTL:     getenvDuration("SESSION_TTL", 7*24*time.Hour),
	}
}

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}

	return fallback
}

func getenvBool(key string, fallback bool) bool {
	raw, ok := os.LookupEnv(key)
	if !ok || raw == "" {
		return fallback
	}

	value, err := strconv.ParseBool(raw)
	if err != nil {
		return fallback
	}

	return value
}

func getenvDuration(key string, fallback time.Duration) time.Duration {
	raw, ok := os.LookupEnv(key)
	if !ok || raw == "" {
		return fallback
	}

	value, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}

	return value
}
