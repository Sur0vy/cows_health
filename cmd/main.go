package main

import (
	"github.com/Sur0vy/cows_health.git/internal/storages"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/Sur0vy/cows_health.git/internal/config"
	"github.com/Sur0vy/cows_health.git/internal/logger"
	"github.com/Sur0vy/cows_health.git/internal/server"
)

func main() {
	cnf := *config.Setup(config.LoadParams())
	log := logger.New(logger.IsDebug(cnf.IsDebug), logger.LogFile(cnf.LogFile))
	log.Info().Msgf("Server start: address: %s", cnf.ServerHostPort)
	defer log.Info().Msg("Server stop")

	log.Info().Msgf("Load storage with parameters %s", cnf.DSN)

	db := connectToDB(cnf.DSN)
	defer db.Close()
	//	createTables(db, log)

	us := storages.NewUserDB(db, log)
	//ds := farm.NewDBStorage(db, log)

	var err = server.SetupServer(us, log).Start(cnf.ServerHostPort)

	if err == nil {
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

/*
func createTables(db *sqlx.DB, log *logger.Logger) {
	ctxIn, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// types
	sqlStr := fmt.Sprint("DO $$ BEGIN " +
		"IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'bolus_type') " +
		"THEN CREATE TYPE bolus_type AS ENUM ('С датчиком PH', 'Без датчика PH'); " +
		"END IF; END$$")
	_, err := db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		log.Panic().Err(err).Msg("Fail then creating type bolus_type")
	}
	log.Info().Msg("Type created: bolus_type")

	//1. user table
	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
		"(%s serial UNIQUE PRIMARY KEY, %s TEXT UNIQUE NOT NULL, %s TEXT NOT NULL)",
		entity.TUser, entity.FUserID, entity.FLogin, entity.FPassword)
	_, err = db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		log.Panic().Err(err).Msgf("Fail then creating table %s", entity.TUser)
	}
	log.Info().Msgf("Table created: %s", entity.TUser)

	//2. breed table
	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
		"(%s serial UNIQUE PRIMARY KEY, %s TEXT NOT NULL)",
		entity.TBreed, entity.FBreedID, entity.FName)
	_, err = db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		log.Panic().Err(err).Msgf("Fail then creating table %s", entity.TBreed)
	}
	log.Info().Msgf("Table created: %s", entity.TBreed)

	//3. farm table
	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
		"(%s serial UNIQUE PRIMARY KEY, %s TEXT NOT NULL, "+
		"%s TEXT UNIQUE NOT NULL, %s INTEGER NOT NULL, "+
		"%s BOOLEAN NOT NULL DEFAULT FALSE)",
		entity.TFarm, entity.FFarmID, entity.FName, entity.FAddress, entity.FUserID, entity.FDeleted)
	_, err = db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		log.Panic().Err(err).Msgf("Fail then creating table %s", entity.TFarm)
	}
	log.Info().Msgf("Table created: %s", entity.TFarm)

	//4. health table
	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
		"(%s INTEGER UNIQUE PRIMARY KEY, %s BOOLEAN, "+
		"%s TEXT, %s TIMESTAMP with time zone, %s BOOLEAN NOT NULL DEFAULT FALSE)",
		entity.THealth, entity.FCowID, entity.FEstrus, entity.FIll, entity.FUpdatedAt, entity.FDeleted)
	_, err = db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		log.Panic().Err(err).Msgf("Fail then creating table %s", entity.THealth)
	}
	log.Info().Msgf("Table created: %s", entity.THealth)

	//5. monitoring data table
	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
		"(%s serial UNIQUE PRIMARY KEY, %s INTEGER NOT NULL, "+
		"%s TIMESTAMP with time zone, %s FLOAT, %s FLOAT, %s FLOAT, %s FLOAT)",
		entity.TMonitoringData, entity.FMDID, entity.FCowID, entity.FAddedAt, entity.FPH, entity.FTemperature, entity.FMovement, entity.FCharge)
	_, err = db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		log.Panic().Err(err).Msgf("Fail then creating table %s", entity.TMonitoringData)
	}
	log.Info().Msgf("Table created: %s", entity.TMonitoringData)

	//6. cow table (проверить, есть ли тип такой)
	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
		"(%s serial UNIQUE PRIMARY KEY, %s TEXT NOT NULL, "+
		"%s INTEGER NOT NULL, %s INTEGER NOT NULL, "+
		"%s INTEGER UNIQUE NOT NULL, %s DATE NOT NULL, "+
		"%s TIMESTAMP with time zone NOT NULL, %s bolus_type NOT NULL, "+
		"%s BOOLEAN NOT NULL DEFAULT FALSE)",
		entity.TCow, entity.FCowID, entity.FName, entity.FBreedID, entity.FFarmID, entity.FBolus, entity.FDateOfBorn, entity.FAddedAt, entity.FBolusType, entity.FDeleted)
	_, err = db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		log.Panic().Err(err).Msgf("Fail then creating table %s", entity.TCow)
	}
	log.Info().Msgf("Table created: %s", entity.TCow)

	//links
	//user-farm link
	sqlStr = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_user_farm; "+
		"ALTER TABLE %s ADD CONSTRAINT fk_user_farm FOREIGN KEY (%s) REFERENCES %s (%s)",
		entity.TFarm, entity.TFarm, entity.FUserID, entity.TUser, entity.FUserID)
	_, err = db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		log.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", entity.TFarm, entity.TUser)
	}
	log.Info().Msgf("foreign key created: %s <-> %s", entity.TFarm, entity.TUser)

	//cow-breed link
	sqlStr = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_cow_breed; "+
		"ALTER TABLE %s ADD CONSTRAINT fk_cow_breed FOREIGN KEY (%s) REFERENCES %s (%s)",
		entity.TCow, entity.TCow, entity.FBreedID, entity.TBreed, entity.FBreedID)
	_, err = db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		log.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", entity.TCow, entity.TBreed)
	}
	log.Info().Msgf("foreign key created: %s <-> %s", entity.TCow, entity.TBreed)

	//cow-farm link
	sqlStr = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_cow_farm; "+
		"ALTER TABLE %s ADD CONSTRAINT fk_cow_farm FOREIGN KEY (%s) REFERENCES %s (%s)",
		entity.TCow, entity.TCow, entity.FFarmID, entity.TFarm, entity.FFarmID)
	_, err = db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		log.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", entity.TCow, entity.TFarm)
	}
	log.Info().Msgf("foreign key created: %s <-> %s", entity.TCow, entity.TFarm)

	//health-cow link
	sqlStr = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_cow_health; "+
		"ALTER TABLE %s ADD CONSTRAINT fk_cow_health FOREIGN KEY (%s) REFERENCES %s (%s)",
		entity.THealth, entity.THealth, entity.FCowID, entity.TCow, entity.FCowID)
	_, err = db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		log.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", entity.THealth, entity.TCow)
	}
	log.Info().Msgf("foreign key created: %s <-> %s", entity.THealth, entity.TCow)

	//monitoring data-cow link
	sqlStr = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_cow_md; "+
		"ALTER TABLE %s ADD CONSTRAINT fk_cow_md FOREIGN KEY (%s) REFERENCES %s (%s)",
		entity.TMonitoringData, entity.TMonitoringData, entity.FCowID, entity.TCow, entity.FCowID)
	_, err = db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		log.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", entity.TMonitoringData, entity.TCow)
	}
	log.Info().Msgf("foreign key created: %s <-> %s", entity.TMonitoringData, entity.TCow)
}
*/
