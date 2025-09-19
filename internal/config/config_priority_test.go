package config

import (
	"os"
	"testing"
)

func TestConfigPriority(t *testing.T) {
	// Create a temporary YAML config
	yamlConfig := &YAMLConfig{
		APIKey:      "yaml_api_key",
		LogLevel:    "debug",
		CacheEnabled: true,
		MatchesPerPage: 5,
	}

	// Test 1: No environment variables - should use YAML values
	os.Clearenv()
	config, err := convertYAMLToConfig(yamlConfig)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if config.FaceitAPIKey != "yaml_api_key" {
		t.Errorf("Expected API key 'yaml_api_key', got '%s'", config.FaceitAPIKey)
	}
	if config.LogLevel != "debug" {
		t.Errorf("Expected log level 'debug', got '%s'", config.LogLevel)
	}
	if !config.CacheEnabled {
		t.Errorf("Expected cache enabled true, got false")
	}
	if config.MatchesPerPage != 5 {
		t.Errorf("Expected matches per page 5, got %d", config.MatchesPerPage)
	}

	// Test 2: Environment variables should override YAML
	os.Setenv("FACEIT_API_KEY", "env_api_key")
	os.Setenv("LOG_LEVEL", "warn")
	os.Setenv("CACHE_ENABLED", "false")
	os.Setenv("MATCHES_PER_PAGE", "15")

	config, err = convertYAMLToConfig(yamlConfig)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if config.FaceitAPIKey != "env_api_key" {
		t.Errorf("Expected API key 'env_api_key', got '%s'", config.FaceitAPIKey)
	}
	if config.LogLevel != "warn" {
		t.Errorf("Expected log level 'warn', got '%s'", config.LogLevel)
	}
	if config.CacheEnabled {
		t.Errorf("Expected cache enabled false, got true")
	}
	if config.MatchesPerPage != 15 {
		t.Errorf("Expected matches per page 15, got %d", config.MatchesPerPage)
	}

	// Test 3: Empty environment variables should fall back to YAML
	os.Setenv("FACEIT_API_KEY", "")
	os.Setenv("LOG_LEVEL", "")
	os.Setenv("CACHE_ENABLED", "")

	config, err = convertYAMLToConfig(yamlConfig)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if config.FaceitAPIKey != "yaml_api_key" {
		t.Errorf("Expected API key 'yaml_api_key', got '%s'", config.FaceitAPIKey)
	}
	if config.LogLevel != "debug" {
		t.Errorf("Expected log level 'debug', got '%s'", config.LogLevel)
	}
	if !config.CacheEnabled {
		t.Errorf("Expected cache enabled true, got false")
	}
}
