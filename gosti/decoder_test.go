package gosti

import (
	"bytes"
	"encoding/json"
	"github.com/renan061/gost"
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

	var (
		// A
		validMapA  = map[string]interface{}{"i": i, "j": j, "k": k}
		validBodyA = &requestBodyMockA{i, j, k}

		// B
		validMapB = map[string]interface{}{
			"bodyA": map[string]interface{}{
				"i": i,
				"j": j,
				"k": k,
			},
			"x": x,
			"y": y,
		}
		validBodyB = &requestBodyMockB{validBodyA, x, y}

		// C
		validMapC = map[string]interface{}{
			"bodyA": map[string]interface{}{
				"i": i,
				"j": j,
				"k": k,
			},
			"bodyB": map[string]interface{}{
				"bodyA": map[string]interface{}{
					"i": i,
					"j": j,
					"k": k,
				},
				"x": x,
				"y": y,
			},
		}
		validBodyC = &requestBodyMockC{validBodyA, validBodyB}
	)

	tests := []struct {
		json           map[string]interface{}
		object         gost.RequestBody
		expectedObject interface{}
		expectedReturn bool
	}{
		{
			// Valid Body A
			json:           validMapA,
			object:         &requestBodyMockA{},
			expectedObject: validBodyA,
			expectedReturn: true,
		}, {
			// Valid Body B
			json:           validMapB,
			object:         &requestBodyMockB{},
			expectedObject: validBodyB,
			expectedReturn: true,
		}, {
			// Valid Body C
			json:           validMapC,
			object:         &requestBodyMockC{},
			expectedObject: validBodyC,
			expectedReturn: true,
		},
	}

	decoder := &BasicDecoder{Responder: &BasicResponder{}}
	var w *httptest.ResponseRecorder
	var r *http.Request
	var jsonBytes []byte
	var ok bool

	for _, test := range tests {
		jsonBytes, _ = json.Marshal(test.json)
		r, _ = http.NewRequest("POST", "/test", bytes.NewReader(jsonBytes))
		w = httptest.NewRecorder()
		ok = decoder.Decode(w, r, test.object)
		r.Body.Close()

		if ok != test.expectedReturn {
			t.Errorf("wrong return: wanted %v, got %v", test.expectedReturn, ok)
		}
		if !reflect.DeepEqual(test.object, test.expectedObject) {
			t.Errorf("unmatching objects: wanted %v, got %v",
				test.expectedObject, test.object)
		}
	}
}
