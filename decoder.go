package gost

import (
	"net/http"
)

var decoder Decoder

type RequestBody interface {
	Valid() (bool, error)
}

type Decoder interface {
	Decode(http.ResponseWriter, *http.Request, RequestBody) bool
}

func GetDecoder() Decoder {
	return decoder
}

func SetDecoder(d Decoder) {
	decoder = d
}
