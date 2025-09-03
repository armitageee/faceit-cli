package config

import (
	"fmt"
	"os"
)

// Config holds the application configuration
type Config struct {
	FaceitAPIKey string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	apiKey := os.Getenv("FACEIT_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("FACEIT_API_KEY environment variable is required")
	}

	return &Config{
		FaceitAPIKey: apiKey,
	}, nil
}
