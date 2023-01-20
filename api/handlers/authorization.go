package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ProSellers/go-honeyapi/internal/database"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

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
		return fiber.ErrForbidden
	}

	fmt.Println(user)

	err = database.Latest.CheckPassword(user, l.Password)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Errorln("User pass error")
		return fiber.ErrForbidden
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

	return ctx.SendString(key)
}
