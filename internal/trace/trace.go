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
	serviceName := os.Getenv("APP_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "ai-passway"
	}
	traceEndpoint := os.Getenv("APP_TRACE_ENDPOINT")
	if traceEndpoint == "" {
		traceEndpoint = "127.0.0.1:4317"
	}

	// resource
	r, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceName(serviceName)))
	if err != nil {
		log.Fatalf("[Trace] create resource failed; %s", err)
	}

	// init ops
	ops := []sdktrace.TracerProviderOption{sdktrace.WithResource(r)}
	if enableTrace {
		e, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(traceEndpoint), otlptracegrpc.WithInsecure())
		if err != nil {
			log.Fatalf("[Trace] create otlp grpc exporter failed; %s", err)
		}
		ops = append(ops, sdktrace.WithBatcher(e))
	}

	// init trace provider
	traceProvider = sdktrace.NewTracerProvider(ops...)
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// init tracer
	tracer = otel.Tracer(serviceName)
}
