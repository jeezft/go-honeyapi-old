package main

import (
	"github.com/ProSellers/go-honeyapi/internal/database"
	"github.com/ProSellers/go-honeyapi/utils/cfg"
)

func main() {
	cfg.Load()
	database.Init()
	startServer()
}
