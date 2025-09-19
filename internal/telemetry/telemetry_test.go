package telemetry

import (
	"context"
	"testing"
)

func TestNew_Disabled(t *testing.T) {
	config := Config{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		Environment:    "test",
		Enabled:        false,
	}

	telemetry, err := New(context.Background(), config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if telemetry == nil {
		t.Fatal("Expected telemetry instance, got nil")
	}

	// Test that tracer works (should be no-op)
	_, span := telemetry.StartSpan(context.Background(), "test.span")
	if span == nil {
		t.Error("Expected span, got nil")
	}
	span.End()

	// Test WithSpan
	err = telemetry.WithSpan(context.Background(), "test.withspan", func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestNew_Enabled_NoExporters(t *testing.T) {
	config := Config{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		Environment:    "test",
		Enabled:        true,
		// No exporters configured
	}

	telemetry, err := New(context.Background(), config)
	if err == nil {
		t.Error("Expected error for no exporters, got nil")
	}

	if telemetry != nil {
		t.Error("Expected nil telemetry for error case")
	}
}

func TestTelemetry_Shutdown(t *testing.T) {
	config := Config{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		Environment:    "test",
		Enabled:        false,
	}

	telemetry, err := New(context.Background(), config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Shutdown should not error
	err = telemetry.Shutdown(context.Background())
	if err != nil {
		t.Errorf("Expected no error on shutdown, got %v", err)
	}
}
