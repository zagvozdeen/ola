package api

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func newViteProxy(viteBase *url.URL) *httputil.ReverseProxy {
	proxy := &httputil.ReverseProxy{
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
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("[proxy] %s %s -> error: %v", r.Method, r.URL.String(), err)
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte("Vite dev server is not reachable"))
		},
	}

	return proxy
}

// rewritePath returns a handler that rewrites r.URL.Path and then proxies.
func rewritePath(proxy *httputil.ReverseProxy, rewrite func(p string) string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = rewrite(r.URL.Path)
		proxy.ServeHTTP(w, r)
	}
}
