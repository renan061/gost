package gosti

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Basic testing needs refactoring into multiple little tests
func TestBasicResponder_InvalidResponse(t *testing.T) {
	responder := &BasicResponder{}
	w := httptest.NewRecorder()
	response := "not_basic_response"
	expectedCode := http.StatusInternalServerError

	responder.Respond(w, response)
	if code := w.Code; code != expectedCode {
		t.Errorf("wrong status code: wanted %v, got %v", expectedCode, code)
	}
}

func TestBasicResponder(t *testing.T) {
	responder := &BasicResponder{}
	var response BasicResponse
	var w *httptest.ResponseRecorder
	var expected string

	// Empty response
	response = BasicResponse{}

	w = httptest.NewRecorder()
	responder.Respond(w, response)
	if code := w.Code; code != http.StatusOK {
		t.Errorf("wrong status code: wanted %v, got %v", http.StatusOK, code)
	}

	// Response with code
	response = BasicResponse{Code: http.StatusCreated}

	w = httptest.NewRecorder()
	responder.Respond(w, response)
	if code := w.Code; code != response.Code {
		t.Errorf("wrong status code: wanted %v, got %v", response.Code, code)
	}

	// Response with message
	response = BasicResponse{Message: "testing the testy thing"}
	expected = `{"message":"` + response.Message + `"}`

	w = httptest.NewRecorder()
	responder.Respond(w, response)
	if body := w.Body.String(); body != expected {
		t.Errorf("wrong body: wanted %v, got %v", expected, body)
	}

	// Response with data
	response = BasicResponse{Message: "message", Data: "data"}
	expected = fmt.Sprintf(`{"message":"%v","data":"%v"}`, response.Message,
		response.Data)

	w = httptest.NewRecorder()
	responder.Respond(w, response)
	if body := w.Body.String(); body != expected {
		t.Errorf("wrong body: wanted %v, got %v", expected, body)
	}
}
