package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Save original environment variables
	originalAPIKey := os.Getenv("FACEIT_API_KEY")
	originalDefaultPlayer := os.Getenv("FACEIT_DEFAULT_PLAYER")
	
	// Clean up after test
	defer func() {
		if originalAPIKey != "" {
			os.Setenv("FACEIT_API_KEY", originalAPIKey)
		} else {
			os.Unsetenv("FACEIT_API_KEY")
		}
		if originalDefaultPlayer != "" {
			os.Setenv("FACEIT_DEFAULT_PLAYER", originalDefaultPlayer)
		} else {
			os.Unsetenv("FACEIT_DEFAULT_PLAYER")
		}
	}()

	tests := []struct {
		name           string
		apiKey         string
		defaultPlayer  string
		expectError    bool
		expectedConfig *Config
	}{
		{
			name:        "valid config with default player",
			apiKey:      "test-api-key",
			defaultPlayer: "testplayer",
			expectError: false,
			expectedConfig: &Config{
				FaceitAPIKey:   "test-api-key",
				DefaultPlayer:  "testplayer",
			},
		},
		{
			name:        "valid config without default player",
			apiKey:      "test-api-key",
			defaultPlayer: "",
			expectError: false,
			expectedConfig: &Config{
				FaceitAPIKey:   "test-api-key",
				DefaultPlayer:  "",
			},
		},
		{
			name:        "missing API key",
			apiKey:      "",
			defaultPlayer: "testplayer",
			expectError: true,
			expectedConfig: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			if tt.apiKey != "" {
				os.Setenv("FACEIT_API_KEY", tt.apiKey)
			} else {
				os.Unsetenv("FACEIT_API_KEY")
			}
			
			if tt.defaultPlayer != "" {
				os.Setenv("FACEIT_DEFAULT_PLAYER", tt.defaultPlayer)
			} else {
				os.Unsetenv("FACEIT_DEFAULT_PLAYER")
			}

			// Test Load function
			config, err := Load()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if config != nil {
					t.Error("Expected nil config but got non-nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if config == nil {
					t.Error("Expected non-nil config but got nil")
					return
				}
				
				if config.FaceitAPIKey != tt.expectedConfig.FaceitAPIKey {
					t.Errorf("FaceitAPIKey = %v, want %v", config.FaceitAPIKey, tt.expectedConfig.FaceitAPIKey)
				}
				if config.DefaultPlayer != tt.expectedConfig.DefaultPlayer {
					t.Errorf("DefaultPlayer = %v, want %v", config.DefaultPlayer, tt.expectedConfig.DefaultPlayer)
				}
			}
		})
	}
}
