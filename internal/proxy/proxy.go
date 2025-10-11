package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/OVINC-CN/AIPassway/internal/logger"
	"github.com/OVINC-CN/AIPassway/internal/utils"
)

func DynamicProxyHandler(w http.ResponseWriter, r *http.Request) {
	// extract service key from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) == 0 || pathParts[0] == "" {
		http.Error(w, "", http.StatusNotImplemented)
		return
	}

	// get real host from environment
	serviceKey := pathParts[0]
	realHostStr := utils.GetRealHostFromEnv(serviceKey)
	if realHostStr == "" {
		http.Error(w, "", http.StatusNotImplemented)
		return
	}

	// parse real host url
	newPath := r.URL.Path[len(serviceKey)+1:]
	newURLStr := strings.TrimSuffix(realHostStr, "/") + newPath
	newURL, err := url.Parse(newURLStr)
	if err != nil {
		logger.Logger().Errorf("parse new url failed: %s\nerror: %v", newURLStr, err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	logger.Logger().Infof("proxying for %s to: %s", r.RemoteAddr, newURL.String())

	// init transport
	transport := &http.Transport{
		// proxy
		Proxy: func(*http.Request) (*url.URL, error) {
			proxyURL := os.Getenv("APP_FORWARD_PROXY_URL")
			if proxyURL == "" {
				return nil, nil
			}
			return url.Parse(proxyURL)
		},
		// close connections after use
		DisableKeepAlives: true,
		// enable compression
		DisableCompression: false,
		// timeout
		IdleConnTimeout:       time.Duration(utils.GetConfigIntFromEnv("APP_IDLE_TIMEOUT", 600)) * time.Second,
		ResponseHeaderTimeout: time.Duration(utils.GetConfigIntFromEnv("APP_HEADER_TIMEOUT", 60)) * time.Second,
	}

	// create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(newURL)
	proxy.Transport = transport

	// modify the request
	originalDirector := proxy.Director
	originalRewrite := proxy.Rewrite
	proxy.Director = nil
	proxy.Rewrite = func(req *httputil.ProxyRequest) {
		// call original rewrite and director if exist
		if originalRewrite != nil {
			originalRewrite(req)
		}
		// call original director if exist
		outReq := req.Out
		if originalDirector != nil {
			originalDirector(outReq)
		}
		// modify the request
		outReq.Host = newURL.Host
		outReq.URL.Path = newPath
		if r.URL.RawQuery != "" {
			outReq.URL.RawQuery = r.URL.RawQuery
		}
		outReq.Header.Del("X-Forwarded-For")
		outReq.Header.Del("X-Real-IP")
	}

	// error handler
	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		logger.Logger().Errorf("proxy error: %v", err)
		http.Error(rw, "", http.StatusBadGateway)
	}

	// serve the request
	proxy.ServeHTTP(w, r)
}
