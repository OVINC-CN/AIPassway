package trace

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer
var traceProvider *sdktrace.TracerProvider

func init() {
	// init ctx
	ctx := context.Background()

	// load config
	enableTrace := os.Getenv("APP_ENABLE_TRACE") != ""
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "ai-passway"
	}
	traceEndpoint := os.Getenv("APP_TRACE_ENDPOINT")
	if traceEndpoint == "" {
		traceEndpoint = "127.0.0.1:4317"
	}

	// resource
	// other use OTEL_RESOURCE_ATTRIBUTES
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName(serviceName)),
	)
	if err != nil {
		log.Fatalf("[Trace] create resource failed; %s", err)
	}

	// init exporter
	var exporter sdktrace.SpanExporter
	if enableTrace {
		var err error
		exporter, err = otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(traceEndpoint), otlptracegrpc.WithInsecure())
		if err != nil {
			log.Fatalf("[Trace] create otlp grpc exporter failed; %s", err)
		}
	} else {
		exporter = newNoopExporter()
	}

	// init trace provider
	traceProvider = sdktrace.NewTracerProvider(
		sdktrace.WithResource(r),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// init tracer
	tracer = otel.Tracer(serviceName)
}
