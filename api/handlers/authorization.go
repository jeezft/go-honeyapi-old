package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ProSellers/go-honeyapi/internal/database"
	"github.com/gofiber/fiber/v2"
)

var ErrLogPass = fiber.NewError(403, "invalid login/password")

type loginendpoint struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Long     bool   `json:"long"`
}

func Authorization(ctx *fiber.Ctx) error {
	var l loginendpoint
	err := json.Unmarshal(ctx.Body(), &l)
	if err != nil {
		return err
	}

	user, err := database.Latest.FindUser(l.Username)
	if err != nil {
		fmt.Println(err)
		return ErrLogPass
	}

	err = database.Latest.CheckPassword(user, l.Password)
	if err != nil {
		// log.WithFields(log.Fields{"err": err}).Errorln("User pass error")
		return ErrLogPass
	}

	ttl := time.Duration(time.Hour * 3)

	if l.Long {
		ttl = time.Duration(time.Hour * 730)
	}

	key, e := database.Latest.CreateSession(user, ttl)
	if e != nil {
		fmt.Println(e)
		return e
	}

	b, err := json.Marshal(&tokenresp{Token: key, UserID: user.ID})
	if err != nil {
		return err
	}

	return ctx.Send(b)
}
