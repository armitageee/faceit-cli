package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// LogLevel represents the logging level
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Service   string                 `json:"service"`
}

// Logger wraps logrus.Logger with Kafka integration
type Logger struct {
	*logrus.Logger
	kafkaWriter *kafka.Writer
	kafkaEnabled bool
	serviceName  string
}

// Config holds logger configuration
type Config struct {
	Level        LogLevel `json:"level"`
	KafkaEnabled bool     `json:"kafka_enabled"`
	KafkaBrokers []string `json:"kafka_brokers"`
	KafkaTopic   string   `json:"kafka_topic"`
	ServiceName  string   `json:"service_name"`
	ProductionMode bool   `json:"production_mode"`
	LogToStdout  bool     `json:"log_to_stdout"`
}

// New creates a new logger instance
func New(config Config) (*Logger, error) {
	// Create logrus logger
	logger := logrus.New()
	
	// Set log level
	level, err := logrus.ParseLevel(string(config.Level))
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}
	logger.SetLevel(level)
	
	// Set JSON formatter for structured logging
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})
	
	// Set output based on configuration
	if config.LogToStdout {
		logger.SetOutput(os.Stdout)
	} else {
		// In production mode without stdout, discard all output
		logger.SetOutput(io.Discard)
	}
	
	l := &Logger{
		Logger:      logger,
		kafkaEnabled: config.KafkaEnabled,
		serviceName:  config.ServiceName,
	}
	
	// Setup Kafka writer if enabled
	if config.KafkaEnabled && len(config.KafkaBrokers) > 0 {
		l.kafkaWriter = &kafka.Writer{
			Addr:     kafka.TCP(config.KafkaBrokers...),
			Topic:    config.KafkaTopic,
			Balancer: &kafka.LeastBytes{},
		}
	}
	
	return l, nil
}

// Log sends a log entry to stdout (if enabled) and Kafka (if enabled)
func (l *Logger) Log(level logrus.Level, message string, fields map[string]interface{}) {
	// Add service name to fields
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["service"] = l.serviceName
	
	// Log to stdout using logrus (only if LogToStdout is enabled)
	entry := l.Logger.WithFields(fields)
	switch level {
	case logrus.DebugLevel:
		entry.Debug(message)
	case logrus.InfoLevel:
		entry.Info(message)
	case logrus.WarnLevel:
		entry.Warn(message)
	case logrus.ErrorLevel:
		entry.Error(message)
	}
	
	// Send to Kafka if enabled
	if l.kafkaEnabled && l.kafkaWriter != nil {
		l.sendToKafka(level, message, fields)
	}
}

// sendToKafka sends log entry to Kafka
func (l *Logger) sendToKafka(level logrus.Level, message string, fields map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level.String(),
		Message:   message,
		Fields:    fields,
		Service:   l.serviceName,
	}
	
	jsonData, err := json.Marshal(entry)
	if err != nil {
		// Fallback to logrus if JSON marshaling fails
		l.Logger.WithError(err).Error("Failed to marshal log entry for Kafka")
		return
	}
	
	// Send to Kafka asynchronously
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		err := l.kafkaWriter.WriteMessages(ctx, kafka.Message{
			Value: jsonData,
		})
		if err != nil {
			// Log error but don't fail the application
			l.Logger.WithError(err).Error("Failed to send log to Kafka")
		}
	}()
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.Log(logrus.DebugLevel, message, f)
}

// Info logs an info message
func (l *Logger) Info(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.Log(logrus.InfoLevel, message, f)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.Log(logrus.WarnLevel, message, f)
}

// Error logs an error message
func (l *Logger) Error(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.Log(logrus.ErrorLevel, message, f)
}

// WithField creates a new logger entry with a field
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

// WithFields creates a new logger entry with fields
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

// Close closes the Kafka writer
func (l *Logger) Close() error {
	if l.kafkaWriter != nil {
		return l.kafkaWriter.Close()
	}
	return nil
}

// ParseLogLevel parses a string log level
func ParseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "warn", "warning":
		return LogLevelWarn
	case "error":
		return LogLevelError
	default:
		return LogLevelInfo
	}
}
