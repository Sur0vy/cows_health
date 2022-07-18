package main

import (
	"github.com/Sur0vy/cows_health.git/config"
	"github.com/Sur0vy/cows_health.git/internal/app"
	"github.com/Sur0vy/cows_health.git/logger"
)

func main() {
	cnf := config.Setup(config.LoadParams())
	log := logger.New(logger.IsDebug(cnf.IsDebug), logger.LogFile(cnf.LogFile))
	log.Info().Msgf("Server start: address: %s", cnf.ServerHostPort)
	defer log.Info().Msg("Server stop")
	log.Info().Msgf("Load storage with parameters %s", cnf.DSN)

	app.Run(cnf, log)
}
