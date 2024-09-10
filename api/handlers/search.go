package handlers

import (
	"encoding/json"
	"sync/atomic"

	"github.com/ProSellers/go-honeyapi/api/wb"
	"github.com/ProSellers/go-honeyapi/internal/database"
	"github.com/ProSellers/go-honeyapi/internal/database/models"
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	Srpx    string `json:"query"`
	Filters `json:"filters"`
}

type Filters struct {
	Page int `json:"page"`
}

type Card struct {
	Name  string
	ID    int
	Price int
	*Info
}

type Info struct {
	CPM    int
	Sales  int
	Stocks []Stocks
}

type Stocks struct {
	Name  string
	Value int
}

type PageInfo struct {
	CardCount  int64 `json:"CardCount"`
	PageNum    int   `json:"Page"`
	HighestCPM int   `json:"HighestCPM"`
}

type SrpxOutput struct {
	Info PageInfo      `json:"Info"`
	Data map[int]*Card `json:"Data"`
}

var ErrTooSmallReq = fiber.NewError(400, "Длина запроса должна превышать 2 символа.")

func Search(ctx *fiber.Ctx) error {
	if _, err := checkauth(ctx); err != nil {
		return err
	}

	var inp Input
	json.Unmarshal(ctx.Body(), &inp)

	if inp.Filters.Page < 1 {
		inp.Filters.Page = 1
	}

	if len(inp.Srpx) < 2 {
		return ErrTooSmallReq
	}

	w := wb.New()

	cards, err := w.Search(inp.Srpx, inp.Filters.Page)
	if err != nil {
		return err
	}

	out := SrpxOutput{Data: make(map[int]*Card)}

	for _, card := range cards {
		i := &Info{
			CPM: card.CPM,
		}

		if card.CPM > out.Info.HighestCPM {
			out.Info.HighestCPM = card.CPM
		}

		c := &Card{
			Name:  card.Name,
			ID:    card.ID,
			Price: card.Price,
			Info:  i,
		}

		atomic.AddInt64(&out.Info.CardCount, 1)

		out.Data[c.ID] = c
	}

	out.Info.PageNum = inp.Page

	// o, err := json.Marshal(out)
	// if err != nil {
	// 	return err
	// }

	return ctx.JSON(out)
}

func checkauth(ctx *fiber.Ctx) (*models.User, error) {
	val := ctx.Request().Header.Peek("authorization")

	if len(val) < 1 {
		return nil, fiber.ErrUnauthorized
	}

	session, err := database.Latest.CheckSession(string(val))
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	u, err := database.Latest.FindUserByID(session.UserID)
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	return u, nil
}
