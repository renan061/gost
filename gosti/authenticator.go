package gosti

import (
	"github.com/renan061/gost"
	"log"
	"net/http"
	"strings"
)

// ==================================================
//
//	JWTAuthenticator
//
// ==================================================

type JWTAuthenticator struct {
	// Responder for errors
	Responder gost.Responder

	// Parses the token
	Parse func(string) (JWTClaims, error)

	// Validates the claims
	Validate func(gost.AuthInfo, JWTClaims) bool

	// Used by the parser
	EncryptionKey []byte
}

type JWTClaims map[string]string

func (a JWTAuthenticator) Authenticate(w http.ResponseWriter, r *http.Request,
	info gost.AuthInfo) (gost.Claims, bool) {

	response := BasicResponse{Code: http.StatusUnauthorized}
	response.Message = "could not authenticate"

	// Checks if authorization header is set
	authStr := r.Header.Get("Authorization")
	if authStr == "" {
		log.Println("gosti.jwtauth: empty header")
		a.Responder.Respond(w, response)
		return nil, false
	}

	// Checks if authorization header is set correctly
	strs := strings.Split(authStr, " ")
	if len(strs) != 2 {
		log.Println("gosti.jwtauth: invalid header")
		a.Responder.Respond(w, response)
		return nil, false
	}
	if scheme := strs[0]; scheme != "Bearer" {
		log.Println("gosti.jwtauth: invalid header")
		a.Responder.Respond(w, response)
		return nil, false
	}

	// Parses the token
	token := strs[1]
	claims, err := a.Parse(token)
	if err != nil {
		log.Println("gosti.jwtauth: invalid token (" + err.Error() + ")")
		a.Responder.Respond(w, response)
		return nil, false
	}

	// Validates the claims
	if !a.Validate(info, claims) {
		log.Println("gosti.jwtauth: could not validate")
		a.Responder.Respond(w, response)
		return nil, false
	}

	return claims, true
}
