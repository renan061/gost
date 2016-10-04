package gost

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Route struct {
	Name        string
	Method      string
	Prefix      string
	Uri         string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

// Defines a decorator as a function that receives a handler and returns
// another handler.
type Decorator func(http.Handler) http.Handler

type Decorators []Decorator

// LogDecorator logs basic information about each request the server receives.
func LogDecorator(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		log.Printf("%s     \t%s\t%s", r.Method, r.RequestURI, time.Since(start))
	})
}

func NewRouter(routes Routes, decorators Decorators) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	var h http.Handler
	for _, route := range routes {
		for _, decorator := range decorators {
			h = decorator(route.HandlerFunc)
		}
		r.Methods(route.Method).Path(route.Prefix + route.Uri).
			Name(route.Name).Handler(h)
	}
	return r
}
