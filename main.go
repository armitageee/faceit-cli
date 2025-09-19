package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"faceit-cli/internal/app"
	"faceit-cli/internal/config"
	"faceit-cli/internal/logger"
	"faceit-cli/internal/telemetry"

	"github.com/joho/godotenv"
)

// Version is set during build time via ldflags
var version = "dev"

func main() {
	// Suppress OpenTelemetry logs immediately to avoid stdout pollution
	// This must be done before any OpenTelemetry initialization
	os.Setenv("OTEL_LOG_LEVEL", "fatal")
	
	// Check for version flag
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Printf("faceit-cli version %s\n", version)
		os.Exit(0)
	}

	// Load environment variables from .env file if it exists
	_ = godotenv.Load() // .env file is optional, so we ignore errors

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	loggerConfig := logger.Config{
		Level:          logger.ParseLogLevel(cfg.LogLevel),
		KafkaEnabled:   cfg.KafkaEnabled,
		KafkaBrokers:   cfg.KafkaBrokers,
		KafkaTopic:     cfg.KafkaTopic,
		ServiceName:    "faceit-cli",
		ProductionMode: cfg.ProductionMode,
		LogToStdout:    cfg.LogToStdout,
	}

	appLogger, err := logger.New(loggerConfig)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Close()

	appLogger.Info("Starting faceit-cli application", map[string]interface{}{
		"version":       version,
		"kafka_enabled": cfg.KafkaEnabled,
		"log_level":     cfg.LogLevel,
	})

	ctx := context.Background()
	
	// Initialize telemetry
	telemetryConfig := telemetry.Config{
		ServiceName:    cfg.ServiceName,
		ServiceVersion: version,
		Environment:    cfg.Environment,
		OTLPEndpoint:   cfg.OTLPEndpoint,
		Enabled:        cfg.TelemetryEnabled,
	}
	
	telemetryInstance, err := telemetry.New(ctx, telemetryConfig)
	if err != nil {
		appLogger.Error("Failed to initialize telemetry", map[string]interface{}{
			"error": err.Error(),
		})
		// Continue without telemetry
		telemetryInstance = &telemetry.Telemetry{}
	}
	defer telemetryInstance.Shutdown(ctx)
	
	application := app.NewApp(cfg, appLogger, telemetryInstance)
	
	if err := application.Run(ctx); err != nil {
		appLogger.Error("Application failed", map[string]interface{}{
			"error": err.Error(),
		})
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	appLogger.Info("Application stopped gracefully")
}
