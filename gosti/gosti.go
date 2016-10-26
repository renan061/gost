package gosti

import (
	"log"
)

const (
	gostiId            = "gosti"
	jwtAuthenticatorId = "jwt_authenticator"
	basicResponderId   = "basic_responder"
)

// Auxiliary function
func logError(id, err string) {
	log.Println(gostiId + "." + id + ": " + err)
}
