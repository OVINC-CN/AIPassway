package trace

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func StartSpan(
	ctx context.Context,
	spanName string,
	spanKind trace.SpanKind,
	options ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	options = append(options, trace.WithSpanKind(spanKind))
	return tracer.Start(ctx, spanName, options...)
}

func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}
