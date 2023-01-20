package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/ProSellers/go-honeyapi/internal/database"
	"github.com/gofiber/fiber/v2"
)

type Userout struct {
	UserID       uint
	Username     string
	SessionID    uint
	CreationDate string
	Roles        int8
}

func Userinfo(ctx *fiber.Ctx) error {
	val := ctx.Request().Header.Peek("authorization")

	if len(val) < 1 {
		fmt.Println("bad header2")
		return fiber.ErrBadRequest
	}

	session, err := database.Latest.CheckSession(string(val))
	if err != nil {
		fmt.Println("cannot verify session")
		return fiber.ErrBadRequest
	}

	fmt.Println("Session:")
	fmt.Println(session)

	u, err := database.Latest.FindUserByID(session.UserID)
	if err != nil {
		fmt.Println("cannot find user")
		return fiber.ErrForbidden
	}

	fmt.Println(u.Username)

	out, err := json.Marshal(&Userout{
		UserID:       u.ID,
		Username:     u.Username,
		SessionID:    session.ID,
		CreationDate: u.CreatedAt.Format("15:04:05-02/01/2006"),
		Roles:        u.Roles,
	})
	if err != nil {
		return fiber.ErrInternalServerError
	}

	ctx.Response().Header.Add("content-type", "application/json")

	return ctx.Send(out)
}
