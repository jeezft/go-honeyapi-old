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

var ErrNoEmail = fiber.NewError(fiber.ErrBadRequest.Code, "Вы не ввели EMAIL")
var ErrNotAnEmail = fiber.NewError(fiber.ErrBadRequest.Code, "EMAIL неверный")

var ErrNoLogin = fiber.NewError(fiber.ErrBadRequest.Code, "Вы не ввели логин")
var ErrLoginIsTooShort = fiber.NewError(fiber.ErrBadRequest.Code, "Логин слишком короткий, минимальная длина 3 символа")

var ErrNoPassword = fiber.NewError(fiber.ErrBadRequest.Code, "Вы не ввели пароль")
var ErrWeakPassword = fiber.NewError(fiber.ErrBadRequest.Code, "Пароль, который вы ввели слишком слаб. Длина пароля должна быть больше 6 символов. Добавьте большие буквы и символы")

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

	if r.Username == "" {
		return ErrNoLogin
	}

	if r.Password == "" {
		return ErrNoPassword
	}

	if r.Email == "" {
		return ErrNoEmail
	}

	if len(r.Username) < 3 {
		return ErrLoginIsTooShort
	}

	if len(r.Password) < 6 || symbolAndCapsRegex.MatchString(r.Password) {
		return ErrWeakPassword
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
