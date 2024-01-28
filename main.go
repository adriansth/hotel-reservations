package main

import (
	"flag"

	"github.com/adriansth/go-hotel-reservations/api"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// server
	listenAddr := flag.String("listenAddr", ":8080", "The listen address of the API server.")
	flag.Parse()
	app := fiber.New()
	apiv1 := app.Group("/api/v1")
	// routes
	apiv1.Get("/user", api.HandleGetUsers)
	apiv1.Get("/user/:id", api.HandleGetUser)
	// server
	app.Listen(*listenAddr)
}
