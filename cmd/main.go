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
	log := logger.New(logger.IsDebug(cnf.IsDebug), logger.LogFile(cnf.LogFile))
	log.Info().Msgf("Server start: address: %s", cnf.ServerHostPort)
	defer log.Info().Msg("Server stop")

	log.Info().Msgf("Load storage with parameters %s", cnf.DSN)

	ds := storage.NewDBStorage(context.Background(), cnf.DSN, log)

	var err = server.SetupServer(ds, ds, log).Start(cnf.ServerHostPort)
	if err == nil {
		log.Panic().Err(err).Msg(err.Error())
	}
}
