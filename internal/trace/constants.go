package trace

import otTrace "go.opentelemetry.io/otel/trace"

const (
	AttributeRequestURI           = "request.uri"
	AttributeRequestRemoteAddr    = "request.remote_addr"
	AttributeRequestContentLength = "request.content_length"
	AttributeStatusCode           = "status.code"
)

const (
	SpanKindInternal = otTrace.SpanKindInternal
	SpanKindServer   = otTrace.SpanKindServer
	SpanKindClient   = otTrace.SpanKindClient
	SpanKindProducer = otTrace.SpanKindProducer
	SpanKindConsumer = otTrace.SpanKindConsumer
)
