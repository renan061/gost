package gost

import (
	"net/http"
)

var authenticator Authenticator

type AuthInfo interface{}

type Claims interface{}

type Authenticator interface {
	Authenticate(http.ResponseWriter, *http.Request, AuthInfo) (Claims, bool)
}

func GetAuthenticator() Authenticator {
	return authenticator
}

func SetAuthenticator(a Authenticator) {
	authenticator = a
}
