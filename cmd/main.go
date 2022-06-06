package main

import (
	"context"
)

func main() {
	cnf := *config.Setup(config.LoadParams())
	ds := storage.NewMapStorage(context.Background(), cnf.DSN)

	err := server.SetupServer(&ds).Run(cnf.ServerHostPort)
	if err == nil {
		panic(err)
	}
}
