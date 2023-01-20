package main

import (
	"github.com/ProSellers/go-honeyapi/api/handlers"
	"github.com/gofiber/fiber/v2"
)

func startServer() error {

	s := fiber.New()

	api := s.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/hello", handlers.HelloWorld)

	//users
	u := v1.Group("/users")
	u.Post("/register", handlers.Register)
	u.Post("/auth", handlers.Authorization)
	u.Get("/whoami", handlers.Userinfo)

	return s.Listen(":8080")
}
