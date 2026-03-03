package config

import (
	"os"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Port           string
	Debug          bool
	DatabasePath   string
	MigrationsPath string
	CORSOrigin     string
	JWTSecret      string
	JWTIssuer      string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		Debug:          getEnv("DEBUG", "false") == "true",
		DatabasePath:   getEnv("DATABASE_PATH", "./data/woo.db"),
		MigrationsPath: getEnv("MIGRATIONS_PATH", "./migrations"),
		CORSOrigin:     getEnv("CORS_ORIGIN", "http://localhost:5173"),
		JWTSecret:      getEnv("JWT_SECRET", "change-me-in-production"),
		JWTIssuer:      getEnv("JWT_ISSUER", "woo-server"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
