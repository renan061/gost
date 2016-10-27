package gosti

import (
	"github.com/renan061/gost"
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
			expectedBody:   "\n", // Does not return ""
			expectedReturn: false,
		}, {
			// Default Response Code
			response:       BasicResponse{},
			expectedCode:   http.StatusOK,
			expectedBody:   "{}",
			expectedReturn: true,
		}, {
			// Given Response Code
			response:       BasicResponse{Code: http.StatusCreated},
			expectedCode:   http.StatusCreated,
			expectedBody:   "{}",
			expectedReturn: true,
		}, {
			// Response with Message
			response:       BasicResponse{Message: msg},
			expectedCode:   http.StatusOK,
			expectedBody:   `{"message":"` + msg + `"}`,
			expectedReturn: true,
		}, {
			// Response with Data
			response:       BasicResponse{Data: data},
			expectedCode:   http.StatusOK,
			expectedBody:   `{"data":"` + data + `"}`,
			expectedReturn: true,
		}, {
			// Response with Message and Data
			response:       BasicResponse{Message: msg, Data: data},
			expectedCode:   http.StatusOK,
			expectedBody:   `{"message":"` + msg + `","data":"` + data + `"}`,
			expectedReturn: true,
		},
	}

	responder := &BasicResponder{}
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
			t.Errorf("wrong body: wanted %v, got %v", test.expectedBody, body)
		}
		if ok != test.expectedReturn {
			t.Errorf("wrong return: wanted %v, got %v", test.expectedReturn, ok)
		}
	}
}
