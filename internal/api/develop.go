package api

import (
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/zagvozdeen/ola/internal/logger"
)

func newViteProxy(log *logger.Logger) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.Out.URL.Scheme = "http"
			r.Out.URL.Host = "localhost:5173"

			r.Out.Header.Set("X-Forwarded-Host", r.In.Host)
			r.Out.Header.Set("X-Forwarded-Proto", "http")

			r.SetXForwarded()
		},
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		ErrorLog: log.GetLog(),
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Errorf("[proxy] %s %s -> error: %v", r.Method, r.URL.String(), err)
			w.WriteHeader(http.StatusBadGateway)
			_, err = w.Write([]byte("Vite dev server is not reachable"))
			if err != nil {
				log.Errorf("[proxy] %s %s -> error: %v", r.Method, r.URL.String(), err)
			}
		},
	}
}
