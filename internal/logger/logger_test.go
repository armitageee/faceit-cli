package logger

import (
	"testing"
)

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected LogLevel
	}{
		{"debug", LogLevelDebug},
		{"DEBUG", LogLevelDebug},
		{"info", LogLevelInfo},
		{"INFO", LogLevelInfo},
		{"warn", LogLevelWarn},
		{"warning", LogLevelWarn},
		{"error", LogLevelError},
		{"ERROR", LogLevelError},
		{"invalid", LogLevelInfo}, // default fallback
		{"", LogLevelInfo},        // default fallback
	}

	for _, test := range tests {
		result := ParseLogLevel(test.input)
		if result != test.expected {
			t.Errorf("ParseLogLevel(%s) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestNewLogger(t *testing.T) {
	config := Config{
		Level:          LogLevelInfo,
		KafkaEnabled:   false,
		ServiceName:    "test-service",
		ProductionMode: false,
		LogToStdout:    true,
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	if logger == nil {
		t.Fatal("Logger is nil")
	}

	if logger.serviceName != "test-service" {
		t.Errorf("Expected service name 'test-service', got '%s'", logger.serviceName)
	}

	if logger.kafkaEnabled != false {
		t.Errorf("Expected kafka disabled, got %v", logger.kafkaEnabled)
	}

	// Test logging methods
	logger.Info("Test info message")
	logger.Warn("Test warning message")
	logger.Error("Test error message")
	logger.Debug("Test debug message") // This might not appear due to log level
}

func TestNewLoggerWithKafka(t *testing.T) {
	config := Config{
		Level:          LogLevelDebug,
		KafkaEnabled:   true,
		KafkaBrokers:   []string{"localhost:9092"},
		KafkaTopic:     "test-topic",
		ServiceName:    "test-service",
		ProductionMode: false,
		LogToStdout:    true,
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create logger with Kafka: %v", err)
	}

	if logger == nil {
		t.Fatal("Logger is nil")
	}

	if logger.kafkaEnabled != true {
		t.Errorf("Expected kafka enabled, got %v", logger.kafkaEnabled)
	}

	if logger.kafkaWriter == nil {
		t.Error("Expected kafka writer to be initialized")
	}

	// Clean up
	logger.Close()
}

func TestNewLoggerProductionMode(t *testing.T) {
	config := Config{
		Level:          LogLevelInfo,
		KafkaEnabled:   false,
		ServiceName:    "test-service",
		ProductionMode: true,
		LogToStdout:    false,
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create production logger: %v", err)
	}

	if logger == nil {
		t.Fatal("Logger is nil")
	}

	// Test logging methods (should not output to stdout)
	logger.Info("Test info message in production mode")
	logger.Warn("Test warning message in production mode")
	logger.Error("Test error message in production mode")
}
