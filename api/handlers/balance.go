package handlers

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/ProSellers/go-honeyapi/internal/database"
	"github.com/ProSellers/go-honeyapi/internal/database/models"
	"github.com/gofiber/fiber/v2"
)

type inp struct {
	Username string
	Value    int
	Action   int
}

// actions
var (
	CASE_ADD = 1
	CASE_SET = 2
)

func SetBalance(ctx *fiber.Ctx) error {
	a, err := checkauth(ctx)
	if err != nil {
		return err
	}

	if err := checkroles(a, models.ROLE_ADMIN); err != nil {
		return fiber.ErrForbidden
	}

	var i inp

	if err := json.Unmarshal(ctx.Body(), &i); err != nil {
		return err
	}

	u, err := database.Latest.FindUser(i.Username)
	if err != nil {
		return errors.New("not found")
	}

	switch i.Action {
	case CASE_ADD:
		{
			u.Balance += i.Value
			database.Latest.UpdateUser(u)
		}
	case CASE_SET:
		{
			u.Balance = i.Value
			database.Latest.UpdateUser(u)
		}
	}

	return ctx.SendString(strconv.Itoa(u.Balance))
}

func GetBalance(ctx *fiber.Ctx) error {
	a, err := checkauth(ctx)
	if err != nil {
		return err
	}

	if err := checkroles(a, models.ROLE_ADMIN); err != nil {
		return fiber.ErrForbidden
	}

	var i inp

	if err := json.Unmarshal(ctx.Body(), &i); err != nil {
		return err
	}

	u, err := database.Latest.FindUser(i.Username)
	if err != nil {
		return errors.New("not found")
	}

	return ctx.SendString(strconv.Itoa(u.Balance))
}

// check for authorization
func checkroles(u *models.User, min int8) error {
	if u.Roles != min {
		return fiber.ErrForbidden
	}

	return nil
}
