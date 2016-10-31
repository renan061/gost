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

type tokenManagerMock struct {
	EncryptionKey []byte
}

func (tm tokenManagerMock) Create(claims JWTClaims) (string, error) {
	token := string(tm.EncryptionKey) + "#"
	for key, value := range claims {
		token = token + key + "," + value + "#"
	}
	return token, nil
}

func (tm tokenManagerMock) Parse(token string) (JWTClaims, error) {
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

func (tm tokenManagerMock) Validate(info gost.AuthInfo, claims JWTClaims) bool {
	return true
}

// ==================================================
//
//	Mocks for testing BasicDecoder
//
// ==================================================

const (
	errRequestBodyMockA = "Body A: missing fields"
	errRequestBodyMockB = "Body B: missing fields"
	errRequestBodyMockC = "Body C: missing fields"
)

type requestBodyMockA struct {
	I string `json:"i"`
	J string `json:"j"`
	K string `json:"k"`
}

type requestBodyMockB struct {
	BodyA *requestBodyMockA `json:"bodyA"`
	X     string            `json:"x"`
	Y     string            `json:"y"`
}

type requestBodyMockC struct {
	BodyA *requestBodyMockA `json:"bodyA"`
	BodyB *requestBodyMockB `json:"bodyB"`
}

type requestBodyMockD struct {
	Z string `json:"z"`
}

func (rb requestBodyMockA) Valid() (bool, error) {
	if rb.I == "" || rb.J == "" || rb.K == "" {
		return false, errors.New(errRequestBodyMockA)
	}
	return true, nil
}

func (rb requestBodyMockB) Valid() (bool, error) {
	if rb.X == "" || rb.Y == "" {
		return false, errors.New(errRequestBodyMockB)
	}
	if rb.BodyA != nil {
		ok, err := rb.BodyA.Valid()
		if !ok {
			return false, err
		}
	}
	return true, nil
}

func (rb requestBodyMockC) Valid() (bool, error) {
	if rb.BodyA == nil || rb.BodyB == nil {
		return false, errors.New(errRequestBodyMockC)
	}
	ok, err := rb.BodyA.Valid()
	if !ok {
		return false, err
	}
	ok, err = rb.BodyB.Valid()
	if !ok {
		return false, err
	}
	return true, nil
}

func (rb requestBodyMockD) Valid() (bool, error) {
	return true, nil
}
