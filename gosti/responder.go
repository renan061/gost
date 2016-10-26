package gosti

import (
	"encoding/json"
	"github.com/renan061/gost"
	"net/http"
)

// ==================================================
//
//	BasicResponder
//
// ==================================================

const basicResponderLogId = "gosti.basic_responder"

type BasicResponse struct {
	PrettyJson bool        `json:"-"`
	Code       int         `json:"-"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

type BasicResponder struct{}

func (_ BasicResponder) Respond(w http.ResponseWriter, r gost.Response) bool {
	// Checks if r is of the correct type
	response, ok := r.(BasicResponse)
	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		logError(basicResponderId, "response type assertion failed")
		return false
	}

	var j []byte
	var err error
	if response.PrettyJson {
		j, err = json.MarshalIndent(response, "", "    ")
	} else {
		j, err = json.Marshal(response)
	}
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		logError(basicResponderId, err.Error())
		return false
	}

	// Default header fields
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Writes the http status code
	if response.Code != 0 {
		w.WriteHeader(response.Code)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	_, err = w.Write(j)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		logError(basicResponderId, err.Error())
		return false
	}

	return true
}
