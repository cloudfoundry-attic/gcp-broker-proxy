package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/urfave/negroni"
)

func ReverseProxy(brokerURL *url.URL) negroni.HandlerFunc {
	reverseProxy := httputil.NewSingleHostReverseProxy(brokerURL)
	dirFunc := reverseProxy.Director

	newDirFunc := func(req *http.Request) {
		dirFunc(req)
		req.Host = brokerURL.Host
	}

	reverseProxy.Director = newDirFunc

	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		reverseProxy.ServeHTTP(rw, r)
		next(rw, r)
	})
}
