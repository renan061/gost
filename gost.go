package gost

type Gost interface {
	Authenticator
	Decoder
	Responder
}
