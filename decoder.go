package gost

import (
	"net/http"
)

type RequestBody interface {
	Valid() (bool, error)
}

type Decoder interface {
	Decode(http.ResponseWriter, *http.Request, RequestBody) bool
}
