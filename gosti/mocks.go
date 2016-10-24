package gosti

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/renan061/gost"
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
	tk := jwt.New(jwt.SigningMethodHS256)
	for key, value := range claims {
		tk.Claims[key] = value
	}
	token, err := tk.SignedString(tm.EncryptionKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (tm tokenManager) Parse(token string) (JWTClaims, error) {
	tk, err := jwt.Parse(token, func(tk *jwt.Token) (interface{}, error) {
		// Checking for signing method (JWT security breach)
		if _, ok := tk.Method.(*jwt.SigningMethodHMAC); !ok {
			str := "unexpected signing method: " + tk.Header["alg"].(string)
			return nil, errors.New(str)
		}
		return tm.EncryptionKey, nil
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

func (tm tokenManager) Validate(info gost.AuthInfo, claims JWTClaims) bool {
	return true
}
