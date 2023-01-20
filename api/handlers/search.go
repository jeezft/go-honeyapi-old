package handlers

import (
	"encoding/json"

	"github.com/ProSellers/go-honeyapi/api/wb"
	"github.com/ProSellers/go-honeyapi/internal/database"
	"github.com/ProSellers/go-honeyapi/internal/database/models"
	"github.com/gofiber/fiber/v2"
)

type Input struct {
	Srpx string `json:"query"`
	Filters
}

type Filters struct {
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

func Search(ctx *fiber.Ctx) error {
	if _, err := checkauth(ctx); err != nil {
		return err
	}

	var inp Input
	json.Unmarshal(ctx.Body(), &inp)

	w := wb.New()

	cards, err := w.Search(inp.Srpx)
	if err != nil {
		return err
	}

	out := make(map[int]*Card)

	for _, card := range cards {
		i := &Info{
			CPM: card.CPM,
		}

		c := &Card{
			Name:  card.Name,
			ID:    card.ID,
			Price: card.Price,
			Info:  i,
		}
		out[c.ID] = c
	}

	o, err := json.Marshal(out)
	if err != nil {
		return err
	}

	return ctx.Send(o)
}

func checkauth(ctx *fiber.Ctx) (*models.User, error) {
	val := ctx.Request().Header.Peek("authorization")

	if len(val) < 1 {
		return nil, fiber.ErrBadRequest
	}

	session, err := database.Latest.CheckSession(string(val))
	if err != nil {
		return nil, fiber.ErrBadRequest
	}

	u, err := database.Latest.FindUserByID(session.UserID)
	if err != nil {
		return nil, fiber.ErrForbidden
	}

	return u, nil
}
