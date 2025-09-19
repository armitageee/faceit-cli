package app

import (
	"context"
	"testing"

	"faceit-cli/internal/config"
	"faceit-cli/internal/logger"
	"faceit-cli/internal/telemetry"
)

// createTestLogger creates a test logger
func createTestLogger() *logger.Logger {
	config := logger.Config{
		Level:          logger.LogLevelInfo,
		KafkaEnabled:   false,
		ServiceName:    "test-service",
		ProductionMode: false,
		LogToStdout:    false, // Disable stdout for tests
	}
	logger, _ := logger.New(config)
	return logger
}

// createTestTelemetry creates a test telemetry instance (disabled)
func createTestTelemetry() *telemetry.Telemetry {
	// Create a proper telemetry instance but with tracing disabled
	ctx := context.Background()
	config := telemetry.Config{
		ServiceName:    "test-service",
		ServiceVersion: "test",
		Environment:    "test",
		Enabled:        false, // Disabled for tests
	}
	
	telemetryInstance, err := telemetry.New(ctx, config)
	if err != nil {
		// If telemetry creation fails, return a nil instance
		// The app should handle nil telemetry gracefully
		return nil
	}
	
	return telemetryInstance
}

func TestNewApp(t *testing.T) {
	tests := []struct {
		name      string
		config    *config.Config
		logger    *logger.Logger
		telemetry *telemetry.Telemetry
		wantErr   bool
	}{
		{
			name: "valid config with cache disabled",
			config: &config.Config{
				FaceitAPIKey: "test-api-key",
				CacheEnabled: false,
			},
			logger:    createTestLogger(),
			telemetry: createTestTelemetry(),
			wantErr:   false,
		},
		{
			name: "valid config with cache enabled",
			config: &config.Config{
				FaceitAPIKey: "test-api-key",
				CacheEnabled: true,
				CacheTTL:     30,
			},
			logger:    createTestLogger(),
			telemetry: createTestTelemetry(),
			wantErr:   false,
		},
		{
			name:      "nil config",
			config:    nil,
			logger:    createTestLogger(),
			telemetry: createTestTelemetry(),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("NewApp() panicked unexpectedly: %v", r)
				}
			}()

			app := NewApp(tt.config, tt.logger, tt.telemetry)
			
			if tt.wantErr {
				if app != nil {
					t.Errorf("NewApp() should have failed but returned app: %v", app)
				}
				return
			}

			if app == nil {
				t.Errorf("NewApp() returned nil app")
				return
			}

			if app.config != tt.config {
				t.Errorf("NewApp() config = %v, want %v", app.config, tt.config)
			}

			if app.logger != tt.logger {
				t.Errorf("NewApp() logger = %v, want %v", app.logger, tt.logger)
			}

			if app.telemetry != tt.telemetry {
				t.Errorf("NewApp() telemetry = %v, want %v", app.telemetry, tt.telemetry)
			}

			if app.repo == nil {
				t.Errorf("NewApp() repo is nil")
			}

			// Test repository type based on cache setting
			if tt.config != nil && tt.config.CacheEnabled {
				// With cache enabled, we expect a cached repository
				// We can't easily test the exact type without exposing internal types
				// So we just verify the repository is not nil
				if app.repo == nil {
					t.Errorf("NewApp() with cache enabled should return a repository")
				}
			} else if tt.config != nil {
				// Without cache, we expect a direct repository
				if app.repo == nil {
					t.Errorf("NewApp() with cache disabled should return a repository")
				}
			}
		})
	}
}

func TestApp_Run(t *testing.T) {
	// This test is limited because we can't easily test the TUI without mocking
	// We'll test the initialization and error handling parts
	
	config := &config.Config{
		FaceitAPIKey: "test-api-key",
		CacheEnabled: false,
	}
	appLogger := createTestLogger()
	
	app := NewApp(config, appLogger, createTestTelemetry())
	
	// Test with a cancelled context to avoid hanging
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately
	
	err := app.Run(ctx)
	
	// The error should be related to context cancellation or TUI initialization
	// We can't easily test the full TUI flow without complex mocking
	if err == nil {
		t.Log("Run() completed without error (this might be expected in test environment)")
	} else {
		t.Logf("Run() returned error (expected in test environment): %v", err)
	}
}

func TestApp_Fields(t *testing.T) {
	config := &config.Config{
		FaceitAPIKey: "test-api-key",
		CacheEnabled: true,
		CacheTTL:     60,
	}
	appLogger := createTestLogger()
	
	app := NewApp(config, appLogger, createTestTelemetry())
	
	// Test that all fields are properly set
	if app.config == nil {
		t.Error("App.config is nil")
	}
	
	if app.repo == nil {
		t.Error("App.repo is nil")
	}
	
	if app.logger == nil {
		t.Error("App.logger is nil")
	}
	
	if app.telemetry == nil {
		t.Error("App.telemetry is nil")
	}
	
	// Test config values
	if app.config.FaceitAPIKey != "test-api-key" {
		t.Errorf("App.config.FaceitAPIKey = %s, want test-api-key", app.config.FaceitAPIKey)
	}
	
	if !app.config.CacheEnabled {
		t.Error("App.config.CacheEnabled should be true")
	}
	
	if app.config.CacheTTL != 60 {
		t.Errorf("App.config.CacheTTL = %d, want 60", app.config.CacheTTL)
	}
}

// Benchmark tests
func BenchmarkNewApp(b *testing.B) {
	config := &config.Config{
		FaceitAPIKey: "test-api-key",
		CacheEnabled: true,
		CacheTTL:     30,
	}
	appLogger := createTestLogger()
	telemetry := createTestTelemetry()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewApp(config, appLogger, telemetry)
	}
}

func BenchmarkNewApp_NoCache(b *testing.B) {
	config := &config.Config{
		FaceitAPIKey: "test-api-key",
		CacheEnabled: false,
	}
	appLogger := createTestLogger()
	telemetry := createTestTelemetry()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewApp(config, appLogger, telemetry)
	}
}
