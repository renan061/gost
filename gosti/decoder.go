package gosti

import (
	"encoding/json"
	"github.com/renan061/gost"
	"io"
	"io/ioutil"
	"net/http"
)

// ==================================================
//
//	BasicDecoder
//
// ==================================================

const (
	ErrDecoderInternal    = "decoder internal error"
	ErrDecoderInvalidBody = "invalid body for request"
)

type basicDecoder struct {
	// Responder for errors
	Responder *basicResponder
}

func NewBasicDecoder() gost.Decoder {
	return &basicDecoder{NewBasicResponder().(*basicResponder)}
}

func (d basicDecoder) Decode(w http.ResponseWriter, r *http.Request,
	rb gost.RequestBody) bool {

	// Reading the request's body
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1*mb))
	if err != nil {
		logError(basicDecoderId, "")
		d.Responder.Respond(w, BasicResponse{
			Message: ErrDecoderInternal,
			Code:    http.StatusInternalServerError,
		})
		return false
	}

	// Parsing the json body
	if err := json.Unmarshal(body, rb); err != nil {
		logError(basicDecoderId, err.Error())
		d.Responder.Respond(w, BasicResponse{
			Message: ErrDecoderInvalidBody,
			Code:    HttpStatusUnprocessableEntity,
		})
		return false
	}

	if rb == nil {
		logError(basicDecoderId, "request body is nil")
		d.Responder.Respond(w, BasicResponse{
			Message: ErrDecoderInvalidBody,
			Code:    HttpStatusUnprocessableEntity,
		})
		return false
	}

	// Checking if the data inside the body is valid
	if ok, err := rb.Valid(); !ok {
		logError(basicDecoderId, err.Error())
		d.Responder.Respond(w, BasicResponse{
			Message: err.Error(),
			Code:    HttpStatusUnprocessableEntity,
		})
		return false
	}

	return true
}
