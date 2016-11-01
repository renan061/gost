package tests

import (
	"errors"
	"github.com/renan061/gost"
	"github.com/renan061/gost/gosti"
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

func (tm tokenManagerMock) Create(claims gosti.JwtClaims) (string, error) {
	token := string(tm.EncryptionKey) + "#"
	for key, value := range claims {
		token = token + key + "," + value + "#"
	}
	return token, nil
}

func (tm tokenManagerMock) Parse(token string) (gosti.JwtClaims, error) {
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

func (tm tokenManagerMock) Validate(info gost.AuthInfo,
	claims gosti.JwtClaims) bool {
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

func (rb requestBodyMockA) Valid() error {
	if rb.I == "" || rb.J == "" || rb.K == "" {
		return errors.New(errRequestBodyMockA)
	}
	return nil
}

func (rb requestBodyMockB) Valid() error {
	if rb.X == "" || rb.Y == "" {
		return errors.New(errRequestBodyMockB)
	}
	if rb.BodyA != nil {
		err := rb.BodyA.Valid()
		if err != nil {
			return err
		}
	}
	return nil
}

func (rb requestBodyMockC) Valid() error {
	if rb.BodyA == nil || rb.BodyB == nil {
		return errors.New(errRequestBodyMockC)
	}
	err := rb.BodyA.Valid()
	if err != nil {
		return err
	}
	err = rb.BodyB.Valid()
	if err != nil {
		return err
	}
	return nil
}

func (rb requestBodyMockD) Valid() error {
	return nil
}
