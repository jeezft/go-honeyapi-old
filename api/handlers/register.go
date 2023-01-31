package handlers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/ProSellers/go-honeyapi/internal/database"
	"github.com/gofiber/fiber/v2"
)

// var client = hcaptcha.New(cfg.Cfg.Captcha.SecretKey)

type reginp struct {
	HcaptchaToken string
	Username      string
	Email         string
	Password      string
}

type tokenresp struct {
	UserID uint
	Token  string
}

var symbolAndCapsRegex = regexp.MustCompile(`(?i)[A-Z][!@#$%^&*()_+-=[]{}\\|;':\"<>,./?\[\]]`)

func Register(ctx *fiber.Ctx) error {
	var r reginp
	err := json.Unmarshal(ctx.Body(), &r)
	if err != nil {
		return err
	}

	if len(r.Username) < 3 || len(r.Password) < 6 || symbolAndCapsRegex.MatchString(r.Password) {
		return fiber.ErrNotAcceptable
	}

	user, e := database.Latest.CreateUser(r.Username, r.Email, r.Password)
	if e != nil {
		return e
	}

	key, e := database.Latest.CreateSession(user, time.Duration(time.Hour*980))
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

// {
//		username: "nigro",
//		password: "cock2288"
// }

//sitekey 30dd0ca3-1d5b-45a6-bde1-26273541219a
//acc secret 0xA6096120b3Ca65B1Fb2D9289378524daa72A22fe
