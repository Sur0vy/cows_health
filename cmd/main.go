package main

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/Sur0vy/cows_health.git/internal/config"
	"github.com/Sur0vy/cows_health.git/internal/logger"
	"github.com/Sur0vy/cows_health.git/internal/server"
	"github.com/Sur0vy/cows_health.git/internal/storages"
)

func main() {
	cnf := *config.Setup(config.LoadParams())
	log := logger.New(logger.IsDebug(cnf.IsDebug), logger.LogFile(cnf.LogFile))
	log.Info().Msgf("Server start: address: %s", cnf.ServerHostPort)
	defer log.Info().Msg("Server stop")

	log.Info().Msgf("Load storage with parameters %s", cnf.DSN)

	db := connectToDB(cnf.DSN)
	defer db.Close()
	log.Info().Msg("Creating database")
	createTables(db)

	us := storages.NewUserDB(db, log)
	fs := storages.NewFarmDB(db, log)
	ms := storages.NewMonitoringDataDB(db, log)
	cs := storages.NewCowDB(db, log)

	var err = server.SetupServer(us, fs, ms, cs, log).Start(cnf.ServerHostPort)

	if err != nil {
		log.Panic().Err(err).Msg(err.Error())
	}
}

func connectToDB(DSN string) *sqlx.DB {
	db, err := sqlx.Open("pgx", DSN)
	if err != nil {
		panic(err)
	}
	return db
}

func createTables(db *sqlx.DB) {
	ctxIn, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db.MustExecContext(ctxIn, DBSchema)
}
