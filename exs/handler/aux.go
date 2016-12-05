package handler

import (
	"errors"
	"github.com/renan061/gost"
	"github.com/renan061/gost/gosti"
	"reflect"
	"strings"
)

func init() {
	tokenManager = &TokenManager{[]byte("encryption_key")}
	server = gosti.NewBasicGost(tokenManager)
}

var (
	server       gost.Gost
	tokenManager *TokenManager
)

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
	// Checks if info and claims are equal (info should be map[string]string)
	return reflect.DeepEqual(info.(map[string]string),
		map[string]string(claims))
}
