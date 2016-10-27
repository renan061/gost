package gosti

import (
	"github.com/renan061/gost"
	"net/http"
)

// ==================================================
//
//	BasicDecoder
//
// ==================================================

type BasicDecoder struct{}

func (_ BasicDecoder) Decode(w http.ResponseWriter, r *http.Request,
	rb gost.RequestBody) bool {

	// // Reading the request's body
	// body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1*mb))
	// if err != nil {
	// 	respond(w, RequestResponse{
	// 		Success: false,
	// 		Message: "internal error",
	// 		code:    http.StatusInternalServerError,
	// 	})
	// 	log.Println(err) // TODO: Logger
	// 	return false
	// }

	// // Closing the reader for the request's body
	// if err := r.Body.Close(); err != nil {
	// 	respond(w, RequestResponse{
	// 		Success: false,
	// 		Message: "internal error",
	// 		code:    http.StatusInternalServerError,
	// 	})
	// 	log.Println(err) // TODO: Logger
	// 	return false
	// }

	// // Parsing the json body
	// if err := json.Unmarshal(body, rb); err != nil {
	// 	respond(w, RequestResponse{
	// 		Success: false,
	// 		Message: "invalid body for request",
	// 		code:    HttpStatusUnprocessableEntity,
	// 	})
	// 	log.Println(err) // TODO: Logger
	// 	return false
	// }

	// if rb.Data == nil {
	// 	respond(w, RequestResponse{
	// 		Success: false,
	// 		Message: "invalid body for request - data must not be null",
	// 		code:    HttpStatusUnprocessableEntity,
	// 	})
	// 	return false
	// }

	// // Checking if the data inside the body is valid
	// if ok, msg := rb.Data.Valid(); !ok {
	// 	respond(w, RequestResponse{
	// 		Success: false,
	// 		Message: msg,
	// 		code:    HttpStatusUnprocessableEntity,
	// 	})
	// 	return false
	// }

	return true
}
