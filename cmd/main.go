package main

import (
	"context"

	"github.com/Sur0vy/cows_health.git/internal/config"
	"github.com/Sur0vy/cows_health.git/internal/server"
	"github.com/Sur0vy/cows_health.git/internal/storage"
)

func main() {
	cnf := *config.Setup(config.LoadParams())
	ds := storage.NewDBStorage(context.Background(), cnf.DSN)

	err := server.SetupServer(&ds).Run(cnf.ServerHostPort)
	if err == nil {
		panic(err)
	}
}
