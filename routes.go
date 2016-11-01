package gost

//
// TODO:
// - Move this from here to gosti
//

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Defines a route decorator as a function that receives a handler and returns
// another handler.
type RouteDecorator func(http.Handler) http.Handler

type RouteDecorators []RouteDecorator

type Route struct {
	Name        string
	Method      string
	Prefix      string
	Uri         string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(routes Routes, decorators RouteDecorators) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	var handler http.Handler
	for _, route := range routes {
		for _, decorator := range decorators {
			handler = decorator(route.HandlerFunc)
		}
		r.Methods(route.Method).Path(route.Prefix + route.Uri).
			Name(route.Name).Handler(handler)
	}
	return r
}
