package gosti

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/renan061/gost"
	"net/http"
	"net/http/httptest"
	"reflect"

	"testing"
)

const (
	jwtBodyErrMsg = `{"message":"could not authenticate"}`
	encryptionKey = "enc_key"
)

var authenticator *JWTAuthenticator

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

func parse(token string) (JWTClaims, error) {
	tk, err := jwt.Parse(token, func(tk *jwt.Token) (interface{}, error) {
		// Checking for signing method (JWT security breach)
		if _, ok := tk.Method.(*jwt.SigningMethodHMAC); !ok {
			str := "unexpected signing method: " + tk.Header["alg"].(string)
			return nil, errors.New(str)
		}
		return []byte(encryptionKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !tk.Valid {
		return nil, errors.New("could not parse token")
	}

	claims := make(map[string]string)
	for key, value := range tk.Claims {
		claims[key] = value.(string)
	}
	return claims, nil
}

func validate(info gost.AuthInfo, claims JWTClaims) bool {
	return true
}

func init() {
	authenticator = &JWTAuthenticator{
		Responder:     &BasicResponder{},
		Validate:      validate,
		Parse:         parse,
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
	testJwt(w, r, t, http.StatusUnauthorized, jwtBodyErrMsg, nil, false)
}

func TestJWTAuthenticator_InvalidHeader(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/test1", nil)
	r.Header.Set("Authorization", "MaybeBearer MaybeToken SomethingElse")
	testJwt(w, r, t, http.StatusUnauthorized, jwtBodyErrMsg, nil, false)

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/test2", nil)
	r.Header.Set("Authorization", "NotBearer Token")
	testJwt(w, r, t, http.StatusUnauthorized, jwtBodyErrMsg, nil, false)
}

func TestJWTAuthenticator_InvalidToken(t *testing.T) {
	token, _ := generateToken(map[string]string{"1": "2"}, "another_enc_key")
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/test", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	testJwt(w, r, t, http.StatusUnauthorized, jwtBodyErrMsg, nil, false)
}

func TestJWTAuthenticator_Ok(t *testing.T) {
	claims := map[string]string{"1": "2", "2": "3", "3": "4"}
	token, _ := generateToken(claims, encryptionKey)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/test", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	testJwt(w, r, t, http.StatusOK, "", claims, true)
}

// Auxiliary
func testJwt(w *httptest.ResponseRecorder, r *http.Request, t *testing.T,
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
			return
		}
	}
}
