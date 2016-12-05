package main

import (
	"github.com/renan061/gost"
	"github.com/renan061/gost/example/handler"
	"github.com/renan061/gost/gosti"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting gost example...")

	r := gost.NewRouter(gost.Routes{
		gost.Route{"Auth", "POST", "", "/auth", handler.AuthHandler},
		gost.Route{"GetUser", "GET", "", "/user", handler.GetUserHandler},
	}, gost.RouteDecorators{
		gosti.LogDecorator,
	})

	log.Fatalln(http.ListenAndServe(":8080", r))
}
