package gosti

import (
	"net/http"
	"net/http/httptest"
	"reflect"

	"testing"
)

const (
	jwtBodyErrMsg = `{"message":"could not authenticate"}`
	jwtEncKey     = "encryption_key"
)

var (
	authenticator *JWTAuthenticator
	tm            *tokenManagerMock
)

func init() {
	tm = &tokenManagerMock{[]byte(jwtEncKey)}
	authenticator = &JWTAuthenticator{
		Responder:    &BasicResponder{},
		TokenManager: tm,
	}
}

// ==================================================
//
//	JWTAuthenticator
//
//	TODO
//	- Refactor to look like responder_test.go
//
// ==================================================

func TestJWTAuthenticator_EmptyHeader(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/test", nil)
	testJwtAuth(w, r, t, http.StatusUnauthorized, jwtBodyErrMsg, nil, false)
}

func TestJWTAuthenticator_InvalidHeader(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/test1", nil)
	r.Header.Set("Authorization", "MaybeBearer MaybeToken SomethingElse")
	testJwtAuth(w, r, t, http.StatusUnauthorized, jwtBodyErrMsg, nil, false)

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/test2", nil)
	r.Header.Set("Authorization", "NotBearer Token")
	testJwtAuth(w, r, t, http.StatusUnauthorized, jwtBodyErrMsg, nil, false)
}

func TestJWTAuthenticator_InvalidToken(t *testing.T) {
	claims := map[string]string{"1": "2"}
	tm.EncryptionKey = []byte("another_encryption_key")
	token, _ := tm.Create(claims)
	tm.EncryptionKey = []byte(jwtEncKey)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/test", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	testJwtAuth(w, r, t, http.StatusUnauthorized, jwtBodyErrMsg, nil, false)
}

func TestJWTAuthenticator_Ok(t *testing.T) {
	claims := map[string]string{"1": "2", "2": "3", "3": "4"}
	token, _ := tm.Create(claims)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/test", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	testJwtAuth(w, r, t, http.StatusOK, "", claims, true)
}

// Auxiliary
func testJwtAuth(w *httptest.ResponseRecorder, r *http.Request, t *testing.T,
	expectedCode int, expectedBody string, expectedClaims map[string]string,
	expectedOk bool) {

	claims, ok := authenticator.Authenticate(w, r, nil)
	if (ok && claims == nil) || (!ok && claims != nil) {
		t.Error("inconsistent return form authenticate")
	}
	if ok != expectedOk {
		t.Errorf("wrong return: wanted %v, got %v", expectedOk, ok)
	}
	if code := w.Code; code != expectedCode {
		t.Errorf("wrong status code: wanted %v, got %v", expectedCode, code)
	}
	if body := w.Body.String(); body != expectedBody {
		t.Errorf("wrong body: wanted %v, got %v", expectedBody, body)
	}
	if claims != nil {
		c := map[string]string(claims.(JWTClaims))
		if !reflect.DeepEqual(expectedClaims, c) {
			t.Errorf("unmatching claims: wanted %v, got %v", expectedClaims, c)
		}
	}
}
