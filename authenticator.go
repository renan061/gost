package gost

import (
	"net/http"
)

type AuthInfo interface{}

type Claims interface{}

type Authenticator interface {
	Authenticate(http.ResponseWriter, *http.Request, AuthInfo) (Claims, bool)
}
