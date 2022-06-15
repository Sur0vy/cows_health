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
		"%s BOOLEAN NOT NULL DEFAULT FALSE)",
		TFarm, FFarmID, FName, FAddress, FUserID, FDeleted)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail then creating table %s", TFarm)
	}
	logger.Wr.Info().Msgf("Table created: %s", TFarm)

	//4. health table
	sqlStr = fmt.Sprintf("CREATE TABLE %s "+
		"(%s INTEGER UNIQUE PRIMARY KEY, %s TEXT, "+
		"%s TEXT, %s TEXT, %s TIMESTAMP, "+
		"%s BOOLEAN NOT NULL DEFAULT FALSE)",
		THealth, FCowID, FDrink, FStress, FIll, FUpdatedAt, FDeleted)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail then creating table %s", THealth)
	}
	logger.Wr.Info().Msgf("Table created: %s", THealth)

	//5. monitoring data table
	sqlStr = fmt.Sprintf("CREATE TABLE %s "+
		"(%s INTEGER PRIMARY KEY AUTOINCREMENT, %s INTEGER NOT NULL, "+
		"%s timestamp, %s FLOAT, %s FLOAT, %s FLOAT, %s FLOAT)",
		TMonitoringData, FMDID, FCowID, FAddedAt, FPH, FTemperature, FMovement, FCharge)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail then creating table %s", TMonitoringData)
	}
	logger.Wr.Info().Msgf("Table created: %s", TMonitoringData)

	//6. cow table (проверить, есть ли тип такой)
	sqlStr = fmt.Sprintf("CREATE TABLE %s "+
		"(%s INTEGER PRIMARY KEY AUTOINCREMENT, %s TEXT NOT NULL, "+
		"%s INTEGER NOT NULL, %s INTEGER NOT NULL, "+
		"%s INTEGER UNIQUE NOT NULL, %s DATE NOT NULL, "+
		"%s TIMESTAMP NOT NULL, %s bolus_type NOT NULL, "+
		"%s BOOLEAN NOT NULL DEFAULT FALSE)",
		TCow, FCowID, FName, FBreedID, FFarmID, FBolus, FDateOfBorn, FAddedAt, FBolusType, FDeleted)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail then creating table %s", TCow)
	}
	logger.Wr.Info().Msgf("Table created: %s", TCow)
	//
	////links
	////user-farm link
	//sqlStr = fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT fk_user_farm FOREIGN KEY (%s) REFERENCES %s (%s)",
	//	TFarm, FUserID, TUser, FUserID)
	//_, err = s.db.ExecContext(ctxIn, sqlStr)
	//if err != nil {
	//	logger.Wr.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", TFarm, TUser)
	//}
	//logger.Wr.Info().Msgf("foreign key created: %s <-> %s", TFarm, TUser)
	//
	////cow-breed link
	//sqlStr = fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT fk_cow_breed FOREIGN KEY (%s) REFERENCES %s (%s)",
	//	TCow, FBreedID, TBreed, FBreedID)
	//_, err = s.db.ExecContext(ctxIn, sqlStr)
	//if err != nil {
	//	logger.Wr.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", TCow, TBreed)
	//}
	//logger.Wr.Info().Msgf("foreign key created: %s <-> %s", TCow, TBreed)
	//
	////cow-farm link
	//sqlStr = fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT fk_cow_farm FOREIGN KEY (%s) REFERENCES %s (%s)",
	//	TCow, FFarmID, TFarm, FFarmID)
	//_, err = s.db.ExecContext(ctxIn, sqlStr)
	//if err != nil {
	//	logger.Wr.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", TCow, TFarm)
	//}
	//logger.Wr.Info().Msgf("foreign key created: %s <-> %s", TCow, TFarm)
	//
	////health-cow link
	//sqlStr = fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT fk_cow_health FOREIGN KEY (%s) REFERENCES %s (%s)",
	//	THealth, FCowID, TCow, FCowID)
	//_, err = s.db.ExecContext(ctxIn, sqlStr)
	//if err != nil {
	//	logger.Wr.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", THealth, TCow)
	//}
	//logger.Wr.Info().Msgf("foreign key created: %s <-> %s", THealth, TCow)
	//
	////monitoring data-cow link
	//sqlStr = fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT fk_cow_md FOREIGN KEY (%s) REFERENCES %s (%s)",
	//	TMonitoringData, FCowID, TCow, FCowID)
	//_, err = s.db.ExecContext(ctxIn, sqlStr)
	//if err != nil {
	//	logger.Wr.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", TMonitoringData, TCow)
	//}
	//logger.Wr.Info().Msgf("foreign key created: %s <-> %s", TMonitoringData, TCow)
}

//TODO внешние ключи
//TODO перечисдение
