package server

import (
	"net/http"

	"github.com/OVINC-CN/AIPassway/internal/logger"
	"github.com/OVINC-CN/AIPassway/internal/middleware"
	"github.com/OVINC-CN/AIPassway/internal/proxy"
)

func Serve() {
	// health check endpoint
	http.HandleFunc(
		"/-/healthz",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		},
	)

	// main proxy handler
	http.Handle(
		"/",
		middleware.PublicAuthMiddleware(middleware.LoggingMiddleware(http.HandlerFunc(proxy.DynamicProxyHandler))),
	)

	// start server
	if err := http.ListenAndServe(":8000", nil); err != nil {
		logger.Logger().Errorf("failed to start server\nerror: %v", err)
	}
}
