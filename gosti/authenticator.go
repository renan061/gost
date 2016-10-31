package gosti

import (
	"github.com/renan061/gost"
	"net/http"
	"strings"
)

// ==================================================
//
//	JWTAuthenticator
//
// ==================================================

type jwtAuthenticator struct {
	// Responder for errors
	Responder *basicResponder

	TokenManager JwtTokenManager
}

type JwtTokenManager interface {
	// Parses the token
	Parse(string) (JwtClaims, error)
	// Validates the claims
	Validate(gost.AuthInfo, JwtClaims) bool
}

type JwtClaims map[string]string

func NewJwtAuthenticator(tokenManager JwtTokenManager) gost.Authenticator {
	return &jwtAuthenticator{
		Responder:    NewBasicResponder().(*basicResponder),
		TokenManager: tokenManager,
	}
}

func (a jwtAuthenticator) Authenticate(w http.ResponseWriter, r *http.Request,
	info gost.AuthInfo) (gost.Claims, bool) {

	response := BasicResponse{Code: http.StatusUnauthorized}
	response.Message = "could not authenticate"

	// Checks if authorization header is set
	authStr := r.Header.Get("Authorization")
	if authStr == "" {
		logError(jwtAuthenticatorId, "empty header")
		a.Responder.Respond(w, response)
		return nil, false
	}

	// Checks if authorization header is set correctly
	strs := strings.Split(authStr, " ")
	if len(strs) != 2 {
		logError(jwtAuthenticatorId, "invalid header")
		a.Responder.Respond(w, response)
		return nil, false
	}
	if scheme := strs[0]; scheme != "Bearer" {
		logError(jwtAuthenticatorId, "invalid header")
		a.Responder.Respond(w, response)
		return nil, false
	}

	// Parses the token
	token := strs[1]
	claims, err := a.TokenManager.Parse(token)
	if err != nil {
		logError(jwtAuthenticatorId, "invalid token ("+err.Error()+")")
		a.Responder.Respond(w, response)
		return nil, false
	}

	// Validates the claims
	if !a.TokenManager.Validate(info, claims) {
		logError(jwtAuthenticatorId, "could not validate")
		a.Responder.Respond(w, response)
		return nil, false
	}

	return claims, true
}
