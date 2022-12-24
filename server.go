package main

import (
	"github.com/ProSellers/go-honeyapi/handlers"
	"github.com/gofiber/fiber/v2"
)

func startServer() error {

	s := fiber.New()

	api := s.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/hello", handlers.HelloWorld)

	return s.Listen(":8080")
}
