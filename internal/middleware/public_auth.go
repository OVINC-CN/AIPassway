package middleware

import (
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/OVINC-CN/AIPassway/internal/logger"
	"github.com/OVINC-CN/AIPassway/internal/trace"
	"github.com/OVINC-CN/AIPassway/internal/utils"
	"github.com/google/uuid"
)

const defaultInternalNetworks = "127.0.0.0/8,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16"

var internalNetworks []*net.IPNet
var passwayAuthToken = os.Getenv("APP_PUBLIC_AUTH_TOKEN")

func init() {
	// get internal networks from environment variable
	envNetworks := os.Getenv("APP_INTERNAL_NETWORKS")
	if envNetworks == "" {
		envNetworks = defaultInternalNetworks
	}

	// parse internal networks from environment variable
	networks := strings.Split(envNetworks, ",")
	for _, network := range networks {
		network = strings.TrimSpace(network)
		if network == "" {
			continue
		}
		_, networkParsed, err := net.ParseCIDR(strings.TrimSpace(network))
		if err != nil {
			logger.Logger().Errorf("failed to parse internal network %s\nerror: %v", network, err)
			os.Exit(1)
		}
		internalNetworks = append(internalNetworks, networkParsed)
		logger.Logger().Infof("added internal network: %s", networkParsed.String())
	}

	// check if auth token is set
	if passwayAuthToken == "" {
		passwayAuthToken = uuid.NewString()
		logger.Logger().Warnf("passway auth token is not set, use random token %s", passwayAuthToken)
	}
}

func isInternalNetwork(clientIP string) bool {
	// no internal networks configured, treat all as external
	if len(internalNetworks) == 0 {
		return false
	}

	// parse client ip
	ip := net.ParseIP(clientIP)
	if ip == nil {
		return false
	}

	// check if ip is in any of the internal networks
	for _, network := range internalNetworks {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

func PublicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// trace
		ctx, span := trace.StartSpan(r.Context(), "PublicAuthMiddleware", trace.SpanKindInternal)
		defer span.End()
		r = r.WithContext(ctx)

		// get client ip address
		clientIP := utils.GetClientIP(r)

		// check if the request is from internal network
		if !isInternalNetwork(clientIP) {
			if authHeader := r.Header.Get("X-AI-Passway-Auth"); authHeader != passwayAuthToken {
				logger.Logger().Warnf("unauthorized access from %s with invalid auth header %s", clientIP, authHeader)
				http.Error(w, "", http.StatusUnauthorized)
				return
			}
			logger.Logger().Infof("external access from %s", clientIP)
		} else {
			logger.Logger().Infof("internal access from %s", clientIP)
		}

		// stop span
		span.End()

		// call the next handler
		next.ServeHTTP(w, r)
	})
}
