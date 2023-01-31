package main

import (
	"github.com/ProSellers/go-honeyapi/internal/database"
	"github.com/ProSellers/go-honeyapi/utils/cfg"
	"github.com/sirupsen/logrus"
)

func main() {
	// wb.GetBrandIdAndUseFirst()ev
	cfg.Load()
	database.Init()
	err := startServer()
	if err != nil {
		logrus.Fatalln(err)
	}
}
