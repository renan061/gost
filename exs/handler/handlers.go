package handler

import (
	"errors"
	"github.com/renan061/gost"
	"github.com/renan061/gost/gosti"
	"log"
	"net/http"
)

// ==================================================
//
//	Auth
//
// ==================================================

type AuthUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u AuthUser) Valid() error {
	if u.Username == "" || u.Password == "" {
		return errors.New("invalid auth body")
	}
	return nil
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	data := &AuthUser{}
	ok := server.Decode(w, r, data)
	if !ok {
		return
	}

	token, err := tokenManager.Create(map[string]string{
		"username": data.Username,
		"password": data.Password,
	})
	if err != nil {
		log.Println(err)
		return
	}

	server.Respond(w, gosti.BasicResponse{
		Data: struct {
			Token string `json:"token"`
		}{token},
	})
}

// ==================================================
//
//	GetUser
//
// ==================================================

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := server.Authenticate(w, r, map[string]string{
		"username": "renan",
		"password": "secret",
	})
	if !ok {
		return
	}

	server.Respond(w, gosti.BasicResponse{
		PrettyJson: true,
		Data: struct {
			A      string      `json:"a"`
			B      string      `json:"b"`
			Claims gost.Claims `json:"claims"`
		}{"A", "B", claims},
	})
}
