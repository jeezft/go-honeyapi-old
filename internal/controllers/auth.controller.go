package controllers

import (
	"github.com/ProSellers/go-honeyapi/internal/database"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	DB *database.Db
}

func NewAuthController(DB *database.Db) AuthController {
	return AuthController{DB}
}

func (ac *AuthController) VerifyCode(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{"succ": "fuck"})
}
