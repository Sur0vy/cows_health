package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Sur0vy/cows_health.git/internal/logger"
)

func NewDBMockStorage(ctx context.Context) *DBStorage {
	s := &DBStorage{}
	s.connectMock()
	s.createMockTables(ctx)
	return s
}

func (s *DBStorage) connectMock() {
	var err error
	s.db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
}

func (s *DBStorage) createMockTables(ctx context.Context) {
	ctxIn, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	//1. user table
	sqlStr := fmt.Sprintf("CREATE TABLE %s "+
		"(%s INTEGER PRIMARY KEY AUTOINCREMENT, %s TEXT UNIQUE NOT NULL, %s TEXT NOT NULL)",
		TUser, FUserID, FLogin, FPassword)
	_, err := s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail then creating table %s", TUser)
	}
	logger.Wr.Info().Msgf("Table created: %s", TUser)

	//2. breed table
	sqlStr = fmt.Sprintf("CREATE TABLE %s "+
		"(%s INTEGER PRIMARY KEY AUTOINCREMENT, %s TEXT NOT NULL)",
		TBreed, FBreedID, FBreed)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail then creating table %s", TBreed)
	}
	logger.Wr.Info().Msgf("Table created: %s", TBreed)

	//3. farm table
	sqlStr = fmt.Sprintf("CREATE TABLE %s "+
		"(%s INTEGER PRIMARY KEY AUTOINCREMENT, %s TEXT NOT NULL, "+
		"%s TEXT UNIQUE NOT NULL, %s INTEGER NOT NULL, "+
		"%s BOOLEAN NOT NULL DEFAULT FALSE, "+
		"FOREIGN KEY (%s) REFERENCES %s(%s))",
		TFarm, FFarmID, FName, FAddress, FUserID, FDeleted,
		FUserID, TUser, FUserID)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail then creating table %s", TFarm)
	}
	logger.Wr.Info().Msgf("Table created: %s", TFarm)

	//4. cow table (проверить, есть ли тип такой)
	sqlStr = fmt.Sprintf("CREATE TABLE %s "+
		"(%s INTEGER PRIMARY KEY AUTOINCREMENT, %s TEXT NOT NULL, "+
		"%s INTEGER NOT NULL, %s INTEGER NOT NULL, "+
		"%s INTEGER UNIQUE NOT NULL, %s DATE NOT NULL, "+
		"%s TIMESTAMP with time zone NOT NULL, %s TEXT NOT NULL, "+
		"%s BOOLEAN NOT NULL DEFAULT FALSE, "+
		"FOREIGN KEY (%s) REFERENCES %s(%s), "+
		"FOREIGN KEY (%s) REFERENCES %s(%s))",
		TCow, FCowID, FName, FBreedID, FFarmID, FBolus, FDateOfBorn, FAddedAt, FBolusType, FDeleted,
		FBreedID, TBreed, FBreedID, FFarmID, TFarm, FFarmID)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail then creating table %s", TCow)
	}
	logger.Wr.Info().Msgf("Table created: %s", TCow)

	//5. health table
	sqlStr = fmt.Sprintf("CREATE TABLE %s "+
		"(%s INTEGER UNIQUE PRIMARY KEY, %s BOOLEAN, "+
		"%s TEXT, %s TIMESTAMP with time zone, "+
		"%s BOOLEAN NOT NULL DEFAULT FALSE, "+
		"FOREIGN KEY (%s) REFERENCES %s(%s))",
		THealth, FCowID, FEstrus, FIll, FUpdatedAt, FDeleted,
		FCowID, TCow, FCowID)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail then creating table %s", THealth)
	}
	logger.Wr.Info().Msgf("Table created: %s", THealth)

	//6. monitoring data table
	sqlStr = fmt.Sprintf("CREATE TABLE %s "+
		"(%s INTEGER PRIMARY KEY AUTOINCREMENT, %s INTEGER NOT NULL, "+
		"%s TIMESTAMP with time zone, %s FLOAT, %s FLOAT, %s FLOAT, %s FLOAT, "+
		"FOREIGN KEY (%s) REFERENCES %s(%s))",
		TMonitoringData, FMDID, FCowID, FAddedAt, FPH, FTemperature, FMovement, FCharge,
		FCowID, TCow, FCowID)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail then creating table %s", TMonitoringData)
	}
	logger.Wr.Info().Msgf("Table created: %s", TMonitoringData)

	//test data
	sqlStr = fmt.Sprintf("INSERT INTO %s(%s) VALUES "+
		"(?), (?), (?)", TBreed, FBreed)
	_, err = s.db.ExecContext(ctxIn, sqlStr, "Голштинская", "Красная датская", "Айрширская")
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail to insert data into %s", TBreed)
	}
	logger.Wr.Info().Msgf("data inserted into %s", TBreed)
}

//TODO внешние ключи
