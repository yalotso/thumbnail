package main

import (
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/go-ozzo/ozzo-routing/slash"
	"github.com/yalotso/thumbnail/api"
	"log"
	"net/http"
)

func main() {
	router := routing.New()

	router.Use(
		// all these handlers are shared by every route
		access.Logger(log.Printf),
		slash.Remover(http.StatusMovedPermanently),
		fault.Recovery(log.Printf),
		content.TypeNegotiator(content.JSON),
	)

	router.Post("/multipart", api.Multipart)
	router.Post("/base64", api.Base64)
	router.Get(`/reference`, api.Reference)

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
