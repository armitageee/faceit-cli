package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds the application configuration
type Config struct {
	FaceitAPIKey      string
	DefaultPlayer     string
	LogLevel          string
	KafkaEnabled      bool
	KafkaBrokers      []string
	KafkaTopic        string
	ProductionMode    bool
	LogToStdout       bool
	MatchesPerPage    int
	MaxMatchesToLoad  int
	CacheEnabled      bool
	CacheTTL          int // Cache TTL in minutes
	ComparisonMatches int // Number of matches to use for comparison
	// Telemetry configuration
	TelemetryEnabled   bool
	OTLPEndpoint       string
	ServiceName        string
	ServiceVersion     string
	Environment        string
}

// Load loads configuration with environment variables taking priority over YAML config
func Load() (*Config, error) {
	var yamlConfig *YAMLConfig
	var err error

	// Try to load YAML config first
	yamlConfig, err = LoadYAMLConfig()
	if err != nil {
		// If YAML config doesn't exist or fails, fall back to environment variables only
		return loadFromEnv()
	}

	// Convert YAML config to Config struct with environment variable overrides
	return convertYAMLToConfig(yamlConfig)
}

// loadFromEnv loads configuration from environment variables (fallback)
func loadFromEnv() (*Config, error) {
	apiKey := os.Getenv("FACEIT_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("FACEIT_API_KEY environment variable is required")
	}

	defaultPlayer := os.Getenv("FACEIT_DEFAULT_PLAYER")
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	kafkaEnabled := os.Getenv("KAFKA_ENABLED") == "true"
	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	if len(kafkaBrokers) == 1 && kafkaBrokers[0] == "" {
		kafkaBrokers = []string{"localhost:9092"}
	}

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		kafkaTopic = "faceit-cli-logs"
	}

	// Parse pagination settings
	matchesPerPage := 10
	if matchesPerPageStr := os.Getenv("MATCHES_PER_PAGE"); matchesPerPageStr != "" {
		if parsed, err := strconv.Atoi(matchesPerPageStr); err == nil && parsed > 0 {
			matchesPerPage = parsed
		}
	}

	maxMatchesToLoad := 100
	if maxMatchesStr := os.Getenv("MAX_MATCHES_TO_LOAD"); maxMatchesStr != "" {
		if parsed, err := strconv.Atoi(maxMatchesStr); err == nil && parsed > 0 {
			maxMatchesToLoad = parsed
		}
	}

	// Parse production mode settings
	productionMode := os.Getenv("PRODUCTION_MODE") == "true"
	logToStdout := os.Getenv("LOG_TO_STDOUT") != "false" // Default to true unless explicitly disabled

	// Parse cache settings
	cacheEnabled := os.Getenv("CACHE_ENABLED") == "true"
	cacheTTL := 30 // Default 30 minutes
	if cacheTTLStr := os.Getenv("CACHE_TTL"); cacheTTLStr != "" {
		if parsed, err := strconv.Atoi(cacheTTLStr); err == nil && parsed > 0 {
			cacheTTL = parsed
		}
	}

	// Parse comparison settings
	comparisonMatches := 20 // Default 20 matches for comparison
	if comparisonStr := os.Getenv("COMPARISON_MATCHES"); comparisonStr != "" {
		if parsed, err := strconv.Atoi(comparisonStr); err == nil && parsed > 0 {
			comparisonMatches = parsed
		}
	}

	// Parse telemetry settings
	telemetryEnabled := os.Getenv("TELEMETRY_ENABLED") == "true"
	otlpEndpoint := os.Getenv("OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "localhost:4317"
	}
	// Zipkin endpoint is handled by OTLP Collector
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "faceit-cli"
	}
	serviceVersion := os.Getenv("SERVICE_VERSION")
	if serviceVersion == "" {
		serviceVersion = "dev"
	}
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}

	return &Config{
		FaceitAPIKey:      apiKey,
		DefaultPlayer:     defaultPlayer,
		LogLevel:          logLevel,
		KafkaEnabled:      kafkaEnabled,
		KafkaBrokers:      kafkaBrokers,
		KafkaTopic:        kafkaTopic,
		ProductionMode:    productionMode,
		LogToStdout:       logToStdout,
		MatchesPerPage:    matchesPerPage,
		MaxMatchesToLoad:  maxMatchesToLoad,
		CacheEnabled:      cacheEnabled,
		CacheTTL:          cacheTTL,
		ComparisonMatches: comparisonMatches,
		TelemetryEnabled:  telemetryEnabled,
		OTLPEndpoint:      otlpEndpoint,
		ServiceName:       serviceName,
		ServiceVersion:    serviceVersion,
		Environment:       environment,
	}, nil
}

// convertYAMLToConfig converts YAMLConfig to Config with environment variable overrides
func convertYAMLToConfig(yamlConfig *YAMLConfig) (*Config, error) {
	// API Key: Environment variable takes priority
	apiKey := os.Getenv("FACEIT_API_KEY")
	if apiKey == "" {
		apiKey = yamlConfig.APIKey
	}
	if apiKey == "" || apiKey == "your_faceit_api_key_here" {
		return nil, fmt.Errorf("API key not configured. Please set FACEIT_API_KEY environment variable or 'api_key' in ~/.config/faceit-cli/config.yml")
	}

	// Helper function to get value with env override
	getStringValue := func(envKey, yamlValue, defaultValue string) string {
		if envValue := os.Getenv(envKey); envValue != "" {
			return envValue
		}
		if yamlValue != "" {
			return yamlValue
		}
		return defaultValue
	}

	getBoolValue := func(envKey string, yamlValue, defaultValue bool) bool {
		if envValue := os.Getenv(envKey); envValue != "" {
			return envValue == "true"
		}
		return yamlValue
	}

	getIntValue := func(envKey string, yamlValue, defaultValue int) int {
		if envValue := os.Getenv(envKey); envValue != "" {
			if parsed, err := strconv.Atoi(envValue); err == nil {
				return parsed
			}
		}
		if yamlValue != 0 {
			return yamlValue
		}
		return defaultValue
	}

	// Parse kafka brokers with env override
	kafkaBrokers := []string{"localhost:9092"}
	if envBrokers := os.Getenv("KAFKA_BROKERS"); envBrokers != "" {
		kafkaBrokers = strings.Split(envBrokers, ",")
	} else if yamlConfig.KafkaBrokers != "" {
		kafkaBrokers = strings.Split(yamlConfig.KafkaBrokers, ",")
	}

	return &Config{
		FaceitAPIKey:      apiKey,
		DefaultPlayer:     getStringValue("FACEIT_DEFAULT_PLAYER", yamlConfig.DefaultPlayer, ""),
		LogLevel:          getStringValue("LOG_LEVEL", yamlConfig.LogLevel, "info"),
		KafkaEnabled:      getBoolValue("KAFKA_ENABLED", yamlConfig.KafkaEnabled, false),
		KafkaBrokers:      kafkaBrokers,
		KafkaTopic:        getStringValue("KAFKA_TOPIC", yamlConfig.KafkaTopic, "faceit-cli-logs"),
		ProductionMode:    getBoolValue("PRODUCTION_MODE", yamlConfig.ProductionMode, false),
		LogToStdout:       getBoolValue("LOG_TO_STDOUT", yamlConfig.LogToStdout, true),
		MatchesPerPage:    getIntValue("MATCHES_PER_PAGE", yamlConfig.MatchesPerPage, 10),
		MaxMatchesToLoad:  getIntValue("MAX_MATCHES_TO_LOAD", yamlConfig.MaxMatchesToLoad, 100),
		CacheEnabled:      getBoolValue("CACHE_ENABLED", yamlConfig.CacheEnabled, false),
		CacheTTL:          getIntValue("CACHE_TTL", yamlConfig.CacheTTL, 30),
		ComparisonMatches: getIntValue("COMPARISON_MATCHES", yamlConfig.ComparisonMatches, 20),
		TelemetryEnabled:  getBoolValue("TELEMETRY_ENABLED", yamlConfig.TelemetryEnabled, false),
		OTLPEndpoint:      getStringValue("OTLP_ENDPOINT", yamlConfig.OTLPEndpoint, "localhost:4317"),
		ServiceName:       getStringValue("SERVICE_NAME", yamlConfig.ServiceName, "faceit-cli"),
		ServiceVersion:    getStringValue("SERVICE_VERSION", yamlConfig.ServiceVersion, "1.0.0"),
		Environment:       getStringValue("ENVIRONMENT", yamlConfig.Environment, "development"),
	}, nil
}
