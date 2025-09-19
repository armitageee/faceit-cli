package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// YAMLConfig represents the YAML configuration structure
type YAMLConfig struct {
	APIKey           string `yaml:"api_key"`
	DefaultPlayer    string `yaml:"default_player"`
	LogLevel         string `yaml:"log_level"`
	KafkaEnabled     bool   `yaml:"kafka_enabled"`
	KafkaBrokers     string `yaml:"kafka_brokers"`
	KafkaTopic       string `yaml:"kafka_topic"`
	ProductionMode   bool   `yaml:"production_mode"`
	LogToStdout      bool   `yaml:"log_to_stdout"`
	MatchesPerPage   int    `yaml:"matches_per_page"`
	MaxMatchesToLoad int    `yaml:"max_matches_to_load"`
	CacheEnabled     bool   `yaml:"cache_enabled"`
	CacheTTL         int    `yaml:"cache_ttl"`
	ComparisonMatches int   `yaml:"comparison_matches"`
	// Telemetry configuration
	TelemetryEnabled bool   `yaml:"telemetry_enabled"`
	OTLPEndpoint     string `yaml:"otlp_endpoint"`
	ServiceName      string `yaml:"service_name"`
	ServiceVersion   string `yaml:"service_version"`
	Environment      string `yaml:"environment"`
	OTELLogLevel     string `yaml:"otel_log_level"`
}

// GetConfigPath returns the appropriate config file path for the current platform
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Use ~/.config/faceit-cli/config.yml
	configDir := filepath.Join(homeDir, ".config", "faceit-cli")
	configFile := filepath.Join(configDir, "config.yml")

	return configFile, nil
}

// LoadYAMLConfig loads configuration from YAML file
func LoadYAMLConfig() (*YAMLConfig, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found at %s", configPath)
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var config YAMLConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	return &config, nil
}

// CreateDefaultConfig creates a default config file with example values
func CreateDefaultConfig() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create default config with all fields
	defaultConfig := YAMLConfig{
		APIKey:           "your_faceit_api_key_here",
		DefaultPlayer:    "",
		LogLevel:         "info",
		KafkaEnabled:     false,
		KafkaBrokers:     "localhost:9092",
		KafkaTopic:       "faceit-cli-logs",
		ProductionMode:   false,
		LogToStdout:      false,
		MatchesPerPage:   10,
		MaxMatchesToLoad: 100,
		CacheEnabled:     true,
		CacheTTL:         30,
		ComparisonMatches: 20,
		TelemetryEnabled: false,
		OTLPEndpoint:     "localhost:4317",
		ServiceName:      "faceit-cli",
		ServiceVersion:   "1.0.0",
		Environment:      "development",
		OTELLogLevel:     "fatal",
	}

	// Marshal to YAML
	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal default config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
