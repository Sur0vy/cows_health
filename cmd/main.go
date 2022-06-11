package main

import (
	"context"

	"github.com/Sur0vy/cows_health.git/internal/config"
	"github.com/Sur0vy/cows_health.git/internal/logger"
	"github.com/Sur0vy/cows_health.git/internal/server"
	"github.com/Sur0vy/cows_health.git/internal/storage"
)

func main() {
	cnf := *config.Setup(config.LoadParams())
	logger.Wr = logger.New(cnf.IsDebug, cnf.LogFile)
	logger.Wr.Info().Msgf("Server start: address:", cnf.ServerHostPort)
	defer logger.Wr.Info().Msg("Server stop")

	logger.Wr.Info().Msgf("Load storage with parameters %s", cnf.DSN)

	ds := storage.NewDBStorage(context.Background(), cnf.DSN)

	err := server.SetupServer(&ds).Run(cnf.ServerHostPort)
	if err == nil {
		logger.Wr.Panic().Err(err).Msg(err.Error())
	}
}
