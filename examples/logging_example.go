package main

import (
	"faceit-cli/internal/logger"
	"os"
	"time"
)

func main() {
	// Check environment variables for production mode
	productionMode := os.Getenv("PRODUCTION_MODE") == "true"
	logToStdout := os.Getenv("LOG_TO_STDOUT") != "false"
	
	// Create logger configuration
	config := logger.Config{
		Level:          logger.LogLevelDebug,
		KafkaEnabled:   false, // Set to true to enable Kafka
		KafkaBrokers:   []string{"localhost:9092"},
		KafkaTopic:     "faceit-cli-logs",
		ServiceName:    "logging-example",
		ProductionMode: productionMode,
		LogToStdout:    logToStdout,
	}

	// Initialize logger
	appLogger, err := logger.New(config)
	if err != nil {
		panic(err)
	}
	defer appLogger.Close()

	// Log different levels
	appLogger.Debug("This is a debug message", map[string]interface{}{
		"component": "example",
		"action":    "startup",
	})

	appLogger.Info("Application started", map[string]interface{}{
		"version": "1.0.0",
		"port":    8080,
	})

	appLogger.Warn("This is a warning message", map[string]interface{}{
		"warning_type": "deprecated_feature",
		"feature":      "old_api",
	})

	appLogger.Error("This is an error message", map[string]interface{}{
		"error_code": "E001",
		"component":  "database",
		"operation":  "connect",
	})

	// Simulate some work
	time.Sleep(100 * time.Millisecond)

	appLogger.Info("Application completed", map[string]interface{}{
		"duration_ms": 100,
		"status":      "success",
	})
}
