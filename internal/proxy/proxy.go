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
	"github.com/sirupsen/logrus"
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
	realHostURL, err := url.Parse(realHostStr)
	if err != nil {
		logger.Logger().Errorf("[DynamicProxyHandler] parse real host url failed: %s\nerror: %v", realHostStr, err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// construct new path without the service key
	newPath := r.URL.Path[len(serviceKey)+1:]

	// init transport
	transport := &http.Transport{
		DisableKeepAlives:     true,
		ResponseHeaderTimeout: time.Duration(utils.GetConfigIntFromEnv("APP_HEADER_TIMEOUT", 60)) * time.Second,
		IdleConnTimeout:       time.Duration(utils.GetConfigIntFromEnv("APP_IDLE_TIMEOUT", 600)) * time.Second,
	}

	// init proxy
	proxyURL := os.Getenv("APP_FORWARD_PROXY_URL")
	if proxyURL != "" {
		// forward proxy url
		forwardProxyURL, err := url.Parse(proxyURL)
		if err != nil {
			logrus.Errorf("[DynamicProxyHandler] parse forward proxy url failed: %s\nerror: %v", proxyURL, err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		// create custom transport with forward proxy
		transport.Proxy = http.ProxyURL(forwardProxyURL)
	}

	// create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(realHostURL)
	proxy.Transport = transport

	// custom director to modify the request
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = realHostURL.Host
		req.URL.Path = newPath
		if r.URL.RawQuery != "" {
			req.URL.RawQuery = r.URL.RawQuery
		}
		req.Header.Set("Connection", "close")
	}

	// add response modifier to ensure connection close
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("Connection", "close")
		return nil
	}

	// serve the request
	proxy.ServeHTTP(w, r)
}
