package handlers

import (
	"encoding/json"
	"regexp"
	"strconv"

	"github.com/ProSellers/go-honeyapi/internal/database"
	"github.com/gofiber/fiber/v2"
)

// var client = hcaptcha.New(cfg.Cfg.Captcha.SecretKey)

type reginp struct {
	HcaptchaToken string
	Username      string
	Password      string
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

	user, e := database.Latest.CreateUser(r.Username, r.Password)
	if e != nil {
		return e
	}

	return ctx.SendString(strconv.Itoa(int(user.ID)))
}

// {
//		username: "nigro",
//		password: "cock2288"
// }

//sitekey 30dd0ca3-1d5b-45a6-bde1-26273541219a
//acc secret 0xA6096120b3Ca65B1Fb2D9289378524daa72A22fe
