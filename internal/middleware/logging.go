package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/OVINC-CN/AIPassway/internal/logger"
	"github.com/OVINC-CN/AIPassway/internal/trace"
	"github.com/OVINC-CN/AIPassway/internal/utils"
	"go.opentelemetry.io/otel/attribute"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// trace
		ctx, span := trace.StartSpan(r.Context(), fmt.Sprintf("%s#%s", r.Method, r.URL.Path), trace.SpanKindServer)
		defer span.End()
		r = r.WithContext(ctx)

		// start timer
		start := time.Now()

		// create a response writer wrapper to capture status code
		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// call the next handler
		next.ServeHTTP(wrappedWriter, r)

		// log response
		duration := time.Since(start)
		logger.Logger().Infof(
			"%s %s %s %s %s %d %s",
			r.RemoteAddr,
			r.Method,
			r.URL.String(),
			utils.FormatBytes(r.ContentLength),
			duration.String(),
			wrappedWriter.statusCode,
			utils.FormatBytes(wrappedWriter.bytesWritten),
		)

		// set span attributes
		span.SetAttributes(
			attribute.String(trace.AttributeRequestURI, r.URL.String()),
			attribute.String(trace.AttributeRequestRemoteAddr, r.RemoteAddr),
			attribute.Int64(trace.AttributeRequestContentLength, r.ContentLength),
			attribute.Int(trace.AttributeStatusCode, wrappedWriter.statusCode),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += int64(n)
	return n, err
}

func (rw *responseWriter) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}
