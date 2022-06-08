package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sur0vy/cows_health.git/internal/logger"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type DBStorage struct {
	db *sql.DB
}

func NewDBStorage(ctx context.Context, DSN string) Storage {
	s := &DBStorage{}
	s.Connect(DSN)
	s.CreateTables(ctx)
	return s
}

func (s *DBStorage) GetFarms(_ context.Context, userID int) (string, error) {
	if userID == 1 {
		var farmsList []Farm
		for i := 1; i < 10; i++ {
			farm := &Farm{
				ID:      i,
				Name:    fmt.Sprintf("Ферма №: %d", i),
				Address: fmt.Sprintf("Адрес фермы №: %d", i),
			}
			farmsList = append(farmsList, *farm)
		}
		if len(farmsList) == 0 {
			return "", errors.New("no farms found")
		}
		data, err := json.Marshal(&farmsList)
		if err != nil {
			return "", err
		}
		return string(data), nil

	} else {
		return "", fmt.Errorf("no farms for user")
	}
}

func (s *DBStorage) AddFarm(ctx context.Context, farm Farm) error {
	return nil
}

func (s *DBStorage) GetFarmInfo(ctx context.Context, farmID int) (string, error) {
	return "", nil
}

func (s *DBStorage) GetCows(_ context.Context, farmID int) (string, error) {
	return "", nil
}

func (s *DBStorage) AddUser(c context.Context, user User) (string, error) {
	return "", nil
}

func (s *DBStorage) CheckUser(_ context.Context, user User) (string, error) {
	return "", nil
}

func (s *DBStorage) GetUser(c context.Context, cookie string) (int, error) {
	return -1, nil
}

func (s *DBStorage) GetBolusesTypes(ctx context.Context) (string, error) {
	return "", nil
}

func (s *DBStorage) GetCowInfo(ctx context.Context, cowID int) (string, error) {
	return "", nil
}

func (s *DBStorage) GetCowBreeds(ctx context.Context) (string, error) {
	return "", nil
}

func (s *DBStorage) AddMonitoringData(ctx context.Context, data MonitoringData) error {
	return nil
}

func (s *DBStorage) DeleteCows(ctx context.Context, IDs []int) error {
	return nil
}

func (s *DBStorage) AddCow(ctx context.Context, cow Cow) error {
	return nil
}

func (s *DBStorage) Connect(DSN string) {
	var err error
	s.db, err = sql.Open("pgx", DSN)
	if err != nil {
		panic(err)
	}
}

func (s *DBStorage) CreateTables(ctx context.Context) {
	ctxIn, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	//1. user table
	sqlStr := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
		"(%s serial UNIQUE PRIMARY KEY, %s TEXT UNIQUE NOT NULL, %s TEXT NOT NULL)",
		TUser, FUserID, FLogin, FPassword)
	_, err := s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail then creating table %s", TUser)
	}
	logger.Wr.Info().Msgf("Table created: %s", TUser)

	//2. breed table
	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
		"(%s serial UNIQUE PRIMARY KEY, %s TEXT NOT NULL)",
		TBreed, FBreedID, FBreed)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail then creating table %s", TBreed)
	}
	logger.Wr.Info().Msgf("Table created: %s", TBreed)

	//3. farm table
	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
		"(%s serial UNIQUE PRIMARY KEY, %s TEXT NOT NULL, "+
		"%s TEXT UNIQUE NOT NULL, %s INTEGER NOT NULL)",
		TFarm, FFarmID, FName, FAddress, FUserID)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail then creating table %s", TFarm)
	}
	logger.Wr.Info().Msgf("Table created: %s", TFarm)

	//4. health table
	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
		"(%s INTEGER UNIQUE PRIMARY KEY, %s TEXT, "+
		"%s TEXT, %s TEXT, %s TIMESTAMP)",
		THealth, FCowID, FDrink, FStress, FIll, FUpdatedAt)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail then creating table %s", THealth)
	}
	logger.Wr.Info().Msgf("Table created: %s", THealth)

	//5. monitoring data table
	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
		"(%s serial UNIQUE PRIMARY KEY, %s INTEGER NOT NULL, "+
		"%s timestamp, %s FLOAT, %s FLOAT, %s FLOAT, %s FLOAT)",
		TMonitoringData, FMDID, FCowID, FAddedAt, FPH, FTemperature, FMovement, FCharge)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail then creating table %s", TMonitoringData)
	}
	logger.Wr.Info().Msgf("Table created: %s", TMonitoringData)

	//6. cow table (проверить, есть ли тип такой)
	sqlStr = fmt.Sprintf("DROP TYPE IF EXISTS bolus_type; "+
		"CREATE TYPE bolus_type AS ENUM ('С датчиком PH', 'Без датчика PH'); "+
		"CREATE TABLE IF NOT EXISTS %s "+
		"(%s serial UNIQUE PRIMARY KEY, %s TEXT NOT NULL, "+
		"%s INTEGER NOT NULL, %s INTEGER NOT NULL, "+
		"%s INTEGER UNIQUE NOT NULL, %s DATE NOT NULL, "+
		"%s TIMESTAMP NOT NULL, %s bolus_type NOT NULL)",
		TCow, FCowID, FName, FBreedID, FFarmID, FBolus, FDateOfBorn, FAddedAt, FBolusType)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail then creating table %s", TCow)
	}
	logger.Wr.Info().Msgf("Table created: %s", TCow)

	//links
	//user-farm link
	sqlStr = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_user_farm; "+
		"ALTER TABLE %s ADD CONSTRAINT fk_user_farm FOREIGN KEY (%s) REFERENCES %s (%s)",
		TFarm, TFarm, FUserID, TUser, FUserID)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", TFarm, TUser)
	}
	logger.Wr.Info().Msgf("foreign key created: %s <-> %s", TFarm, TUser)

	//cow-breed link
	sqlStr = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_cow_breed; "+
		"ALTER TABLE %s ADD CONSTRAINT fk_cow_breed FOREIGN KEY (%s) REFERENCES %s (%s)",
		TCow, TCow, FBreedID, TBreed, FBreedID)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", TCow, TBreed)
	}
	logger.Wr.Info().Msgf("foreign key created: %s <-> %s", TCow, TBreed)

	//cow-farm link
	sqlStr = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_cow_farm; "+
		"ALTER TABLE %s ADD CONSTRAINT fk_cow_farm FOREIGN KEY (%s) REFERENCES %s (%s)",
		TCow, TCow, FFarmID, TFarm, FFarmID)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", TCow, TFarm)
	}
	logger.Wr.Info().Msgf("foreign key created: %s <-> %s", TCow, TFarm)

	//health-cow link
	sqlStr = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_cow_health; "+
		"ALTER TABLE %s ADD CONSTRAINT fk_cow_health FOREIGN KEY (%s) REFERENCES %s (%s)",
		THealth, THealth, FCowID, TCow, FCowID)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", THealth, TCow)
	}
	logger.Wr.Info().Msgf("foreign key created: %s <-> %s", THealth, TCow)

	//monitoring data-cow link
	sqlStr = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_cow_md; "+
		"ALTER TABLE %s ADD CONSTRAINT fk_cow_md FOREIGN KEY (%s) REFERENCES %s (%s)",
		TMonitoringData, TMonitoringData, FCowID, TCow, FCowID)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", TMonitoringData, TCow)
	}
	logger.Wr.Info().Msgf("foreign key created: %s <-> %s", TMonitoringData, TCow)
}

func (s *DBStorage) Ping() error {
	return s.db.Ping()
}
