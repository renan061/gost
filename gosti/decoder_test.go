package gosti

import (
	"bytes"
	"encoding/json"
	"github.com/renan061/gost"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ==================================================
//
//	BasicDecoder
//
// ==================================================

func TestBasicDecoder(t *testing.T) {
	// var (
	// 	validBodyA   = requestBodyMockA{I: "i_value", J: "j_value", K: "k_value"}
	// 	invalidBodyA = requestBodyMockA{I: "i_value", J: "j_value", K: ""}
	// )

	tests := []struct {
		json           string
		object         gost.RequestBody
		expectedObject interface{}
		expectedReturn bool
	}{
		// ioutil.ReadAll
		// json.Unmarshal
		// rb == nil
		// rb.Valid()
		{
			json:           `{"a":"b"}`,
			object:         &requestBodyMockA{},
			expectedObject: nil, // TODO: Rename
			expectedReturn: false,
		}, {
			json:           "{}",
			object:         &requestBodyMockA{},
			expectedObject: nil, // TODO: Rename
			expectedReturn: false,
		},
	}

	decoder := &BasicDecoder{}
	var w *httptest.ResponseRecorder
	var r *http.Request
	var ok bool

	for _, test := range tests {
		j, _ := json.Marshal(test.json)
		r, _ = http.NewRequest("POST", "/test", bytes.NewBuffer(j))
		w = httptest.NewRecorder()
		ok = decoder.Decode(w, r, test.object)

		if ok != test.expectedReturn {
			t.Errorf("wrong return: wanted %v, got %v", test.expectedReturn, ok)
		}
	}
}
