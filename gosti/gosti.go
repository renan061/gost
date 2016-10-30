package gosti

import (
	"log"
)

const (
	gostiId            = "gosti"
	jwtAuthenticatorId = "jwt_authenticator"
	basicResponderId   = "basic_responder"
	basicDecoderId     = "basic_decoder"

	mb = 1048576 // 1048576 bytes == 1 megabyte

	HttpStatusUnprocessableEntity = 422
)

// Auxiliary function
func logError(id, err string) {
	log.Println(gostiId + "." + id + ": " + err)
}
