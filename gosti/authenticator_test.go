package gosti

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/renan061/gost"
	"net/http"
	"net/http/httptest"

	"testing"
)

const (
	jwtErrMsg     = `{"message":"could not authenticate"}`
	encryptionKey = "encryption_key"
)

var authenticator *JWTAuthenticator

func init() {
	authenticator = &JWTAuthenticator{
		Responder:     &BasicResponder{},
		Validator:     &JWTValidatorMock{},
		EncryptionKey: []byte(encryptionKey),
	}
}

// ==================================================
//
//	JWTAuthenticator
//
// ==================================================

func TestJWTAuthenticator_EmptyHeader(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/test", nil)
	testJwt(w, r, t, nil, http.StatusUnauthorized, jwtErrMsg)
}

func TestJWTAuthenticator_InvalidHeader(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/test1", nil)
	r.Header.Set("Authorization", "MaybeBearer MaybeToken SomethingElse")
	testJwt(w, r, t, nil, http.StatusUnauthorized, jwtErrMsg)

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/test2", nil)
	r.Header.Set("Authorization", "NotBearer Token")
	testJwt(w, r, t, nil, http.StatusUnauthorized, jwtErrMsg)
}

func TestJWTAuthenticator_InvalidToken(t *testing.T) {
	claims := map[string]string{"1": "2", "2": "3", "3": "4"}
	token, _ := generateToken(claims, "another_encryption_key")
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/test", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	testJwt(w, r, t, claims, http.StatusUnauthorized, jwtErrMsg)
}

func TestJWTAuthenticator_Ok(t *testing.T) {
	claims := map[string]string{"1": "2", "2": "3", "3": "4"}
	token, _ := generateToken(claims, encryptionKey)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/test", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	testJwt(w, r, t, claims, http.StatusOK, "")
}

// ==================================================
//
//	Auxiliary
//
// ==================================================

func testJwt(w *httptest.ResponseRecorder, r *http.Request, t *testing.T,
	expectedClaims gost.AuthInfo, expectedCode int, expectedBody string) {

	authenticator.Authenticate(w, r, expectedClaims)
	if code := w.Code; code != expectedCode {
		t.Errorf("wrong status code: wanted %v, got %v", expectedCode, code)
	}
	if body := w.Body.String(); body != expectedBody {
		t.Errorf("wrong body: wanted %v, got %v", expectedBody, body)
	}
}

type JWTValidatorMock struct{}

func (v JWTValidatorMock) Validate(info gost.AuthInfo, claims JWTClaims) bool {
	// expectedClaims := info.(map[string]string)
	// if reflect.DeepEqual(claims, expectedClaims) {
	// 	return false
	// }
	return true
}

// Generic generate token function
func generateToken(claims JWTClaims, encryptionKey string) (string, error) {
	tk := jwt.New(jwt.SigningMethodHS256)
	for key, value := range claims {
		tk.Claims[key] = value
	}
	token, err := tk.SignedString([]byte(encryptionKey))
	if err != nil {
		return "", err
	}
	return token, nil
}
