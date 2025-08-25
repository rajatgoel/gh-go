package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config holds application configuration
type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME" default:"gh-go-frontend"`
	Environment string `envconfig:"ENVIRONMENT" default:"development"`
	Port        int    `envconfig:"PORT" default:"5051"`
}

// Load loads configuration from environment variables and .env file
func Load() (*Config, error) {
	// Load .env file if it exists, but don't fail if it doesn't
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}

	return &config, nil
}
