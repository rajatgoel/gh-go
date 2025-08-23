package frontend

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// Config holds OpenTelemetry configuration
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
}

// DefaultConfig returns default OTEL configuration
func DefaultConfig() *Config {
	return &Config{
		ServiceName:    "gh-go-frontend",
		ServiceVersion: "1.0.0",
		Environment:    "development",
	}
}

// SetupOTEL initializes OpenTelemetry with the given configuration
func SetupOTEL(ctx context.Context, config *Config) (func(), error) {
	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(config.ServiceName),
			semconv.ServiceVersion(config.ServiceVersion),
			semconv.DeploymentEnvironment(config.Environment),
		),
	)
	if err != nil {
		return nil, err
	}

	// Setup trace exporter
	traceExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return nil, err
	}

	// Setup trace provider
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(traceProvider)

	// Setup propagation
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Setup metrics - prometheus exporter provides its own meter provider
	_, err = prometheus.New()
	if err != nil {
		return nil, err
	}

	// Create cleanup function
	cleanup := func() {
		if err := traceProvider.Shutdown(ctx); err != nil {
			slog.Error("failed to shutdown trace provider", "error", err)
		}
	}

	return cleanup, nil
}
