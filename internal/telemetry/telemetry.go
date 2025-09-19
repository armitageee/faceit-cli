package telemetry

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// Config holds telemetry configuration
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	OTLPEndpoint   string
	Enabled        bool
}

// Telemetry manages OpenTelemetry tracing
type Telemetry struct {
	tracerProvider *sdktrace.TracerProvider
	tracer         trace.Tracer
	shutdown       func(context.Context) error
}

// NewDisabled creates a disabled telemetry instance for testing
func NewDisabled() *Telemetry {
	return &Telemetry{
		tracer: noop.NewTracerProvider().Tracer("faceit-cli-disabled"),
	}
}

// New creates a new telemetry instance
func New(ctx context.Context, cfg Config) (*Telemetry, error) {
	if !cfg.Enabled {
		// Return a no-op telemetry instance
		return &Telemetry{
			tracer: noop.NewTracerProvider().Tracer("faceit-cli"),
		}, nil
	}

	// OTLP logs are suppressed by setting OTEL_LOG_LEVEL in main.go

	// Create resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironment(cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create exporters
	var exporters []sdktrace.SpanExporter

	// OTLP gRPC exporter
	if cfg.OTLPEndpoint != "" {
		// For gRPC, we need to remove http:// prefix and use just host:port
		endpoint := strings.TrimPrefix(cfg.OTLPEndpoint, "http://")
		endpoint = strings.TrimPrefix(endpoint, "https://")
		
		// Using OTLP gRPC endpoint
		
		otlpExporter, err := otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(endpoint),
			otlptracegrpc.WithInsecure(), // For development
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
		}
		exporters = append(exporters, otlpExporter)
	}

	// Note: Zipkin export is handled by OTLP Collector
	// Direct Zipkin export is removed to use proper OTLP → Collector → Zipkin flow

	if len(exporters) == 0 {
		return nil, fmt.Errorf("no exporters configured")
	}

	// Create tracer provider with first exporter
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporters[0]),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Add additional exporters if any
	for i := 1; i < len(exporters); i++ {
		tp.RegisterSpanProcessor(sdktrace.NewBatchSpanProcessor(exporters[i]))
	}

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Set global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &Telemetry{
		tracerProvider: tp,
		tracer:         tp.Tracer("faceit-cli"),
		shutdown:       tp.Shutdown,
	}, nil
}

// Tracer returns the tracer instance
func (t *Telemetry) Tracer() trace.Tracer {
	return t.tracer
}

// Shutdown gracefully shuts down the telemetry
func (t *Telemetry) Shutdown(ctx context.Context) error {
	if t.shutdown != nil {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return t.shutdown(ctx)
	}
	return nil
}

// StartSpan creates a new span with the given name and options
func (t *Telemetry) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name, opts...)
}

// WithSpan executes a function within a span
func (t *Telemetry) WithSpan(ctx context.Context, name string, fn func(context.Context) error, opts ...trace.SpanStartOption) error {
	ctx, span := t.StartSpan(ctx, name, opts...)
	defer span.End()

	return fn(ctx)
}
