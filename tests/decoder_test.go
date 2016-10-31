package tests

import (
	"bytes"
	"encoding/json"
	"github.com/renan061/gost"
	"github.com/renan061/gost/gosti"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// ==================================================
//
//	BasicDecoder
//
// ==================================================

func TestBasicDecoder(t *testing.T) {
	const (
		i = "i_value"
		j = "j_value"
		k = "k_value"
		x = "x_value"
		y = "y_value"
	)

	// Look for requestBodyMock structs in mocks.go
	var (
		// A
		validMapA    = map[string]interface{}{"i": i, "j": j, "k": k}
		validBodyA   = &requestBodyMockA{i, j, k}
		invalidMapA  = map[string]interface{}{"i": i, "j": j, "k": ""}
		invalidBodyA = &requestBodyMockA{i, j, ""}

		// B
		validMapB   = map[string]interface{}{"x": x, "y": y}
		validBodyB  = &requestBodyMockB{X: x, Y: y}
		invalidMapB = map[string]interface{}{"bodyA": invalidMapA,
			"x": x, "y": y}
		invalidBodyB = &requestBodyMockB{invalidBodyA, x, y}

		// C
		validMapC = map[string]interface{}{"bodyA": validMapA,
			"bodyB": validMapB}
		validBodyC   = &requestBodyMockC{validBodyA, validBodyB}
		invalidMapC  = map[string]interface{}{"bodyB": invalidMapB}
		invalidBodyC = &requestBodyMockC{BodyB: invalidBodyB}
	)

	tests := []struct {
		json           interface{}
		object         gost.RequestBody
		expectedObject interface{}
		expectedReturn bool
		expectedCode   int
		expectedBody   string
	}{
		{
			// 1: Valid A
			json:           validMapA,
			object:         &requestBodyMockA{},
			expectedObject: validBodyA,
			expectedReturn: true,
			expectedCode:   http.StatusOK,
			expectedBody:   "",
		}, {
			// 2: Invalid A
			json:           invalidMapA,
			object:         &requestBodyMockA{},
			expectedObject: invalidBodyA,
			expectedReturn: false,
			expectedCode:   gosti.HttpStatusUnprocessableEntity,
			expectedBody:   newResponseBody(errRequestBodyMockA, ""),
		}, {
			// 3: Valid B
			json:           validMapB,
			object:         &requestBodyMockB{},
			expectedObject: validBodyB,
			expectedReturn: true,
			expectedCode:   http.StatusOK,
			expectedBody:   "",
		}, {
			// 4: Invalid B
			json:           invalidMapB,
			object:         &requestBodyMockB{},
			expectedObject: invalidBodyB,
			expectedReturn: false,
			expectedCode:   gosti.HttpStatusUnprocessableEntity,
			expectedBody:   newResponseBody(errRequestBodyMockA, ""),
		}, {
			// 5: Valid C
			json:           validMapC,
			object:         &requestBodyMockC{},
			expectedObject: validBodyC,
			expectedReturn: true,
			expectedCode:   http.StatusOK,
			expectedBody:   "",
		}, {
			// 6: Invalid C
			json:           invalidMapC,
			object:         &requestBodyMockC{},
			expectedObject: invalidBodyC,
			expectedReturn: false,
			expectedCode:   gosti.HttpStatusUnprocessableEntity,
			expectedBody:   newResponseBody(errRequestBodyMockC, ""),
		}, {
			// 7: Empty body
			json:           map[string]interface{}{"irrelevant": "data"},
			object:         &requestBodyMockD{},
			expectedObject: &requestBodyMockD{},
			expectedReturn: true,
			expectedCode:   http.StatusOK,
			expectedBody:   "",
		}, {
			// 7: JSON null
			json:           new(string),
			object:         &requestBodyMockD{},
			expectedObject: &requestBodyMockD{},
			expectedReturn: false,
			expectedCode:   gosti.HttpStatusUnprocessableEntity,
			expectedBody:   newResponseBody(gosti.ErrDecoderInvalidBody, ""),
		},
	}

	decoder := gosti.NewBasicDecoder()
	var w *httptest.ResponseRecorder
	var r *http.Request
	var jsonBytes []byte
	var ok bool

	for num, test := range tests {
		jsonBytes, _ = json.Marshal(test.json)
		r, _ = http.NewRequest("POST", "/test", bytes.NewReader(jsonBytes))
		w = httptest.NewRecorder()
		ok = decoder.Decode(w, r, test.object)
		r.Body.Close()

		if !reflect.DeepEqual(test.object, test.expectedObject) {
			t.Errorf("%v: unmatching objects: wanted %v, got %v", num+1,
				test.expectedObject, test.object)
		}
		if ok != test.expectedReturn {
			t.Errorf("%v: wrong return: wanted %v, got %v", num+1,
				test.expectedReturn, ok)
		}
		if code := w.Code; code != test.expectedCode {
			t.Errorf("%v: wrong status code: wanted %v, got %v", num+1,
				test.expectedCode, code)
		}
		if body := w.Body.String(); body != test.expectedBody {
			t.Errorf("%v: wrong body: wanted %v, got %v", num+1,
				test.expectedBody, body)
		}
	}
}
