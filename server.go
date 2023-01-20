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

	//wb
	w := v1.Group("/wb")
	w.Post("/search", handlers.Search)
	w.Post("/addtoken", handlers.AddToken)
	w.Get("/tokens", handlers.GetTokens)

	//admin
	a := v1.Group("/admin")
	// a.Post("/getuser")
	a.Post("/balance", handlers.SetBalance).
		Get("/balance", handlers.GetBalance)

	return s.Listen(":8080")
}
