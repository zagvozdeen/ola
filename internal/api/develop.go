package api

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/zagvozdeen/ola/internal/logger"
)

func newViteProxy(log *logger.Logger, viteBase *url.URL) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Rewrite: func(preq *httputil.ProxyRequest) {
			preq.Out.URL.Scheme = "http"
			preq.Out.URL.Host = viteBase.Host

			preq.Out.Header.Set("X-Forwarded-Host", preq.In.Host)
			preq.Out.Header.Set("X-Forwarded-Proto", "http")

			preq.SetXForwarded()
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
