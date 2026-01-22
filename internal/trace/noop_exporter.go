package trace

import (
	"context"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// noopExporter is a SpanExporter that discards all spans.
type noopExporter struct{}

// newNoopExporter creates a new noopExporter.
func newNoopExporter() sdktrace.SpanExporter {
	return &noopExporter{}
}

// ExportSpans discards all spans without any processing.
func (e *noopExporter) ExportSpans(_ context.Context, _ []sdktrace.ReadOnlySpan) error {
	return nil
}

// Shutdown performs no operation and returns nil.
func (e *noopExporter) Shutdown(_ context.Context) error {
	return nil
}
