package gost

import (
	"net/http"
)

type Response interface{}

type Responder interface {
	Respond(http.ResponseWriter, Response) bool
}
