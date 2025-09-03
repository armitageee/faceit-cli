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
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
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

	return &Config{
		FaceitAPIKey:     apiKey,
		DefaultPlayer:    defaultPlayer,
		LogLevel:         logLevel,
		KafkaEnabled:     kafkaEnabled,
		KafkaBrokers:     kafkaBrokers,
		KafkaTopic:       kafkaTopic,
		ProductionMode:   productionMode,
		LogToStdout:      logToStdout,
		MatchesPerPage:   matchesPerPage,
		MaxMatchesToLoad: maxMatchesToLoad,
		CacheEnabled:     cacheEnabled,
		CacheTTL:         cacheTTL,
	}, nil
}
