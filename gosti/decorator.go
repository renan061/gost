package gosti

import (
	"log"
	"net/http"
	"time"
)

// LogDecorator logs basic information (method, uri, elapsed time) about each
// request the server receives.
func LogDecorator(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		str := "%s     \t%s\t%s"
		log.Printf(str, r.Method, r.RequestURI, time.Since(start))
	})
}
