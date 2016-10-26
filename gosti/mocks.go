package gosti

import (
	"errors"
	"github.com/renan061/gost"
	"strings"
)

// ==================================================
//
//	Mocks for testing JWTAuthenticator
//
// ==================================================

type tokenManager struct {
	EncryptionKey []byte
}

func (tm tokenManager) Create(claims JWTClaims) (string, error) {
	token := string(tm.EncryptionKey) + "#"
	for key, value := range claims {
		token = token + key + "," + value + "#"
	}
	return token, nil
}

func (tm tokenManager) Parse(token string) (JWTClaims, error) {
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

func (tm tokenManager) Validate(info gost.AuthInfo, claims JWTClaims) bool {
	return true
}
