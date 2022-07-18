package app

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/Sur0vy/cows_health.git/config"
	"github.com/Sur0vy/cows_health.git/internal/server"
	"github.com/Sur0vy/cows_health.git/internal/storages"
	"github.com/Sur0vy/cows_health.git/logger"
	"github.com/Sur0vy/cows_health.git/migrations"
)

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
	db.MustExecContext(ctxIn, migrations.DBSchema)
}

func Run(cnf *config.Config, log *logger.Logger) {
	db := connectToDB(cnf.DSN)
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Panic().Err(err).Msg(err.Error())
		}
	}(db)
	log.Info().Msg("Creating database")
	createTables(db)

	repo := storages.NewStorageDB(db, log)

	var err = server.SetupServer(*repo, log).Start(cnf.ServerHostPort)

	if err != nil {
		log.Panic().Err(err).Msg(err.Error())
	}
}
