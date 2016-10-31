package main

import (
	"errors"
	"github.com/renan061/gost"
	"github.com/renan061/gost/gosti"
	"log"
	"net/http"
	"strings"
)

var server *gost.Gost

func main() {
	r := gost.NewRouter(gost.Routes{
		gost.Route{"Test 1", "GET", "/test", "/foo", FooHandler},
		gost.Route{"Test 2", "GET", "/test", "/bar", BarHandler},
	}, gost.RouteDecorators{
		gosti.LogDecorator,
	})

	encryptionKey := []byte("encryption_key")
	authenticator := gosti.NewJwtAuthenticator(&TokenManager{encryptionKey})
	decoder := gosti.NewBasicDecoder()
	responder := gosti.NewBasicResponder()

	server = &gost.Gost{authenticator, decoder, responder}

	log.Fatalln(http.ListenAndServe(":8080", r))
}

// ==================================================
//
//	Handlers
//
// ==================================================

func FooHandler(w http.ResponseWriter, r *http.Request) {
	server.Respond(w, gosti.BasicResponse{PrettyJson: true, Data: "data"})
}

func BarHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("bar")
}

// ==================================================
//
//	Token Manager
//
// ==================================================

type TokenManager struct {
	EncryptionKey []byte
}

func (tm TokenManager) Create(claims gosti.JwtClaims) (string, error) {
	token := string(tm.EncryptionKey) + "#"
	for key, value := range claims {
		token = token + key + "," + value + "#"
	}
	return token, nil
}

func (tm TokenManager) Parse(token string) (gosti.JwtClaims, error) {
	arr := strings.Split(token, "#")
	if arr[0] != string(tm.EncryptionKey) {
		return nil, errors.New("could not parse token")
	}

	arr = arr[1 : len(arr)-1]
	claims := make(map[string]string)
	for _, e := range arr {
		keyvalue := strings.Split(e, ",")
		claims[keyvalue[0]] = keyvalue[1]
	}
	return claims, nil
}

func (tm TokenManager) Validate(info gost.AuthInfo,
	claims gosti.JwtClaims) bool {
	return true
}
