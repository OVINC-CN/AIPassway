package middleware

import (
	"net/http"
	"time"

	"github.com/OVINC-CN/AIPassway/internal/logger"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// create a response writer wrapper to capture status code
		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// call the next handler
		next.ServeHTTP(wrappedWriter, r)

		// log response
		duration := time.Since(start)
		logger.Logger().Infof(
			"[LoggingMiddleware] %s %s %s %dB %dms %d %dB",
			r.RemoteAddr,
			r.Method,
			r.URL.String(),
			r.ContentLength,
			duration.Milliseconds(),
			wrappedWriter.statusCode,
			wrappedWriter.bytesWritten,
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
