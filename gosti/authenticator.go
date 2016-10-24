package gosti

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
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

type JWTValidator interface {
	Validate(gost.AuthInfo, JWTClaims) bool
}

type JWTClaims map[string]string

type JWTAuthenticator struct {
	// Responder for errors
	Responder gost.Responder

	// Validates the claims
	Validator JWTValidator

	// Used by the parser
	EncryptionKey []byte
}

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
	claims, err := a.parse(token)
	if err != nil {
		log.Println("gosti.jwtauth: invalid token (" + err.Error() + ")")
		a.Responder.Respond(w, response)
		return nil, false
	}

	// Validates the claims
	if !a.Validator.Validate(info, claims) {
		log.Println("gosti.jwtauth: could not validate")
		a.Responder.Respond(w, response)
		return nil, false
	}

	return claims, true
}

// TODO: Remove this and require it to be an external dependency
func (a JWTAuthenticator) parse(token string) (JWTClaims, error) {
	tk, err := jwt.Parse(token, func(tk *jwt.Token) (interface{}, error) {
		// Checking for signing method (JWT security breach)
		if _, ok := tk.Method.(*jwt.SigningMethodHMAC); !ok {
			str := "unexpected signing method: " + tk.Header["alg"].(string)
			return nil, errors.New(str)
		}
		return a.EncryptionKey, nil
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
