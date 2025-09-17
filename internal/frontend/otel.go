package frontend

import (
	"context"
	"log/slog"
	"time"

	"github.com/earthboundkid/versioninfo/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	"github.com/rajatgoel/gh-go/internal/config"
)

// SetupOTEL initializes OpenTelemetry with the given configuration and returns a cleanup function
// that handles shutdown automatically with proper timeout and error handling.
func SetupOTEL(ctx context.Context, cfg *config.Config) (func(), error) {
	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(versioninfo.Version),
			semconv.DeploymentEnvironment(cfg.Environment),
		),
	)
	if err != nil {
		return nil, err
	}

	// Setup trace provider
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(traceProvider)

	// Setup metrics provider
	metricProvider := metric.NewMeterProvider(
		metric.WithResource(res),
	)
	otel.SetMeterProvider(metricProvider)

	// Setup propagation
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Return cleanup function that handles context and timeout internally
	cleanup := func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := traceProvider.Shutdown(shutdownCtx); err != nil {
			slog.Error("failed to shutdown trace provider", "error", err)
		}
		if err := metricProvider.Shutdown(shutdownCtx); err != nil {
			slog.Error("failed to shutdown metric provider", "error", err)
		}
	}

	return cleanup, nil
}
