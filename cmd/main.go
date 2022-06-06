package main

import (
	"cmd/main.go/internal/server"
	"context"

	"github.com/Sur0vy/cows_health.git/internal/config"
)

func main() {
	cnf := *config.Setup(config.LoadParams())
	ds := storage.NewMapStorage(context.Background(), cnf.DSN)

	err := server.SetupServer(&ds).Run(cnf.ServerHostPort)
	if err == nil {
		panic(err)
	}
}
