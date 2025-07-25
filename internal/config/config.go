package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config holds all configuration for the application, loaded from environment variables.
type Config struct {
	// Server
	Port string `envconfig:"PORT" default:"3001"`

	// Redis
	RedisAddr     string `envconfig:"REDIS_ADDR" required:"true"`
	RedisPassword string `envconfig:"REDIS_PASSWORD"`

	// Auth
	JWTSecret       string `envconfig:"JWT_SECRET" required:"true"`
	JWTExpiresHours int    `envconfig:"JWT_EXPIRES_HOURS" default:"24"`

	// Logging
	LogLevel  string `envconfig:"LOG_LEVEL" default:"info"`
	LogFormat string `envconfig:"LOG_FORMAT" default:"text"`
}

// Load reads configuration from environment variables.
// It will first attempt to load a .env file if present.
func Load() (*Config, error) {
	// Attempt to load .env file, useful for local development.
	// In production, environment variables should be set directly.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	var cfg Config
	err := envconfig.Process("", &cfg)
	return &cfg, err
}