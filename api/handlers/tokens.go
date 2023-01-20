package handlers

import (
	"encoding/json"

	"github.com/ProSellers/go-honeyapi/internal/database"
	"github.com/gofiber/fiber/v2"
)

type tinp struct {
	Token string
}

func GetTokens(ctx *fiber.Ctx) error {
	u, err := checkauth(ctx)
	if err != nil {
		return fiber.ErrForbidden
	}

	tokens, err := database.Latest.GetTokens(u)
	if err != nil {
		return fiber.ErrNotFound
	}

	o, err := json.Marshal(&tokens)
	if err != nil {
		return err
	}

	return ctx.Send(o)
}

func AddToken(ctx *fiber.Ctx) error {
	u, err := checkauth(ctx)
	if err != nil {
		return fiber.ErrForbidden
	}

	var i tinp
	err = json.Unmarshal(ctx.Body(), &i)
	if err != nil {
		return err
	}

	t, err := database.Latest.AddToken(u, i.Token)
	if err != nil {
		return err
	}

	o, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return ctx.Send(o)
}
