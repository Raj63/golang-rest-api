package tracer

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/Raj63/golang-rest-api/pkg/infrastructure/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials"
)

type (

	// Span wraps up trace.Span which is the individual component of a trace. It represents a single
	// named and timed operation of a workflow that is traced. A Tracer is used to create a Span and
	// it is then up to the operation the Span represents to properly end the Span when the operation
	// itself ends.
	Span struct {
		trace.Span
	}

	// Tracer wraps up trace.Tracer which is the creator of Spans.
	Tracer struct {
		trace.Tracer
	}

	// Option is wraps up trace.Option and applies an option to a TracerConfig.
	Option struct {
		trace.TracerOption
	}

	// Provider provides access to instrumentation Tracers.
	Provider struct {
		*sdktrace.TracerProvider
	}

	// ProviderConfig indicates how the OpenTelemetry tracer provider should be initialised.
	ProviderConfig struct {
		// ServiceName indicates the service name to trace.
		ServiceName string

		// CollectorAddress indicates the collector's address.
		CollectorAddress string

		// CollectorBasicAuthCreds indicates the collector's basic auth credentials.
		CollectorBasicAuthCreds string

		// CollectorConnectTimeout indicate the duration to timeout when connecting to the collector.
		// By default, it is 5 * time.Second.
		CollectorConnectTimeout time.Duration
	}
)

// NewTracer initialises the OpenTelemetry tracer.
func NewTracer(config config.AppConfig) (Tracer, *Provider, func() error, error) {
	tracerProvider, traceProviderShutdownHandler, err := NewTracerProvider(
		&ProviderConfig{
			ServiceName:      config.ServiceName,
			CollectorAddress: config.TracerCollectorAddress,
		},
	)
	if err != nil {
		return Tracer{}, nil, nil, err
	}

	return newTracer(config.ServiceName), tracerProvider, traceProviderShutdownHandler, err
}

// NewTracerProvider initializes the provider for OpenTelemetry traces.
func NewTracerProvider(c *ProviderConfig) (*Provider, func() error, error) {
	ctx := context.Background()
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(c.ServiceName),
		),
	)
	if err != nil {
		return nil, nil, err
	}

	defaultTracerProviderConfig(c)

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(c.CollectorAddress),
		otlptracegrpc.WithTimeout(c.CollectorConnectTimeout),
	}

	if c.CollectorBasicAuthCreds != "" {
		opts = append(opts,
			otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
			otlptracegrpc.WithHeaders(
				map[string]string{
					"Authorization": fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(c.CollectorBasicAuthCreds))),
				},
			),
		)
	}

	// Initialise the tracer provider which will use the exporter to push traces to the collector.
	exporter, err := otlptracegrpc.New(
		ctx,
		opts...,
	)

	if err != nil {
		return nil, nil, err
	}

	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(exporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(batchSpanProcessor),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return &Provider{tracerProvider}, func() error {
		if err := exporter.Shutdown(ctx); err != nil {
			return err
		}

		if err := tracerProvider.Shutdown(ctx); err != nil {
			return err
		}

		return nil
	}, nil
}

// newTracer initialises a tracer.
func newTracer(instrumentationName string, opts ...trace.TracerOption) Tracer {
	provider := otel.GetTracerProvider()

	return Tracer{
		provider.Tracer(instrumentationName, opts...),
	}
}

// SpanFromContext returns the current Span from ctx.
//
// If no Span is currently set in ctx an implementation of a Span that performs no operations is
// returned.
func SpanFromContext(ctx context.Context) trace.Span {
	return Span{
		trace.SpanFromContext(ctx),
	}
}

// GetTraceIDFromContext return the trace ID from the context
func GetTraceIDFromContext(ctx context.Context) string {
	return SpanFromContext(ctx).SpanContext().TraceID().String()
}

func defaultTracerProviderConfig(c *ProviderConfig) {
	if c.CollectorConnectTimeout == 0 {
		c.CollectorConnectTimeout = 5 * time.Second
	}
}
