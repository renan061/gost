package gost

import (
	"net/http"
)

var responder Responder

type Response interface{}

type Responder interface {
	Respond(http.ResponseWriter, Response) bool
}

func GetResponder() Responder {
	return responder
}

func SetResponder(r Responder) {
	responder = r
}
