package tests

import (
	"github.com/renan061/gost"
	"github.com/renan061/gost/gosti"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ==================================================
//
//	BasicResponder
//
//	TODO
//	- Test PrettyJson option
//	- Test json.Marshal error
//	- Test w.Write error
//
// ==================================================

func TestBasicResponder(t *testing.T) {
	const (
		msg  = "testing the testy test"
		data = "data many datas"
	)

	tests := []struct {
		response       gost.Response
		expectedCode   int
		expectedBody   string
		expectedReturn bool
	}{
		{
			// Invalid response
			response:       "not_basic_response",
			expectedCode:   http.StatusInternalServerError,
			expectedBody:   "\n", // Does not return empty string
			expectedReturn: false,
		}, {
			// Default Response Code
			response:       gosti.BasicResponse{},
			expectedCode:   http.StatusOK,
			expectedReturn: true,
		}, {
			// Given Response Code
			response:       gosti.BasicResponse{Code: http.StatusCreated},
			expectedCode:   http.StatusCreated,
			expectedReturn: true,
		}, {
			// Response with Message
			response:       gosti.BasicResponse{Message: msg},
			expectedCode:   http.StatusOK,
			expectedBody:   newResponseBody(msg, ""),
			expectedReturn: true,
		}, {
			// Response with Data
			response:       gosti.BasicResponse{Data: data},
			expectedCode:   http.StatusOK,
			expectedBody:   newResponseBody("", data),
			expectedReturn: true,
		}, {
			// Response with Message and Data
			response:       gosti.BasicResponse{Message: msg, Data: data},
			expectedCode:   http.StatusOK,
			expectedBody:   newResponseBody(msg, data),
			expectedReturn: true,
		},
	}

	responder := gosti.NewBasicResponder()
	var w *httptest.ResponseRecorder
	var ok bool

	for _, test := range tests {
		w = httptest.NewRecorder()
		ok = responder.Respond(w, test.response)

		if code := w.Code; code != test.expectedCode {
			t.Errorf("wrong status code: wanted %v, got %v",
				test.expectedCode, code)
		}
		if body := w.Body.String(); body != test.expectedBody {
			t.Errorf("wrong body: wanted %x, got %x", test.expectedBody, body)
		}
		if ok != test.expectedReturn {
			t.Errorf("wrong return: wanted %v, got %v", test.expectedReturn, ok)
		}
	}
}
