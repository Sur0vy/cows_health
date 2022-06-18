package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/mattn/go-sqlite3"

	"github.com/Sur0vy/cows_health.git/internal/logger"
)

type DBStorage struct {
	db *sql.DB
}

func NewDBStorage(ctx context.Context, DSN string) *DBStorage {
	s := &DBStorage{}
	s.connect(DSN)
	s.createTables(ctx)
	return s
}

func (s *DBStorage) AddUser(c context.Context, user User) (string, error) {
	userHash := getMD5Hash(user.Login)
	u := s.GetUser(c, userHash)
	if u != nil {
		logger.Wr.Warn().Msgf("User %v already exists", user.Login)
		return userHash, NewExistError(fmt.Sprintf("user %s already exist", user.Login))
	}

	passwordHash, err := getCryptoPassword(user.Password)
	if err != nil {
		logger.Wr.Warn().Msg("Error than encrypting password")
		return "", err
	}

	ctxIn, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	//добавление
	sqlStr := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES ($1, $2)",
		TUser, FLogin, FPassword)
	_, err = s.db.ExecContext(ctxIn, sqlStr, userHash, passwordHash)
	if err != nil {
		logger.Wr.Warn().Err(err).Msg("db request error")
		return "", err
	}
	return userHash, nil
}

func (s *DBStorage) GetUser(ctx context.Context, userHash string) *User {
	ctxIn, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user := &User{}

	sqlStr := fmt.Sprintf("SELECT %s, %s FROM %s WHERE %s = $1",
		FUserID, FPassword, TUser, FLogin)
	row := s.db.QueryRowContext(ctxIn, sqlStr, userHash)
	err := row.Scan(&user.ID, &user.Password)

	if err != nil {
		logger.Wr.Warn().Err(err).Msg("db request error")
		return nil
	}
	return user
}

func (s *DBStorage) GetUserHash(c context.Context, user User) (string, error) {
	userHash := getMD5Hash(user.Login)
	u := s.GetUser(c, userHash)
	if u == nil {
		logger.Wr.Warn().Msgf("User %v not exists", user.Login)
		return "", NewEmptyError(fmt.Sprintf("user %s not exists", user.Login))
	}

	if checkPassword(u.Password, user.Password) {
		return userHash, nil
	}
	return "", NewEmptyError(fmt.Sprintf("password wrong for user %s", user.Login))
}

func (s *DBStorage) GetFarms(c context.Context, userID int) (string, error) {
	ctxIn, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var farms []Farm
	sqlStr := fmt.Sprintf("SELECT %s, %s, %s FROM %s WHERE %s = $1 AND NOT %s",
		FFarmID, FName, FAddress, TFarm, FUserID, FDeleted)
	rows, err := s.db.QueryContext(ctxIn, sqlStr, userID)

	if err != nil {
		logger.Wr.Warn().Err(err).Msg("db request error")
		return "", err
	}

	defer func() {
		_ = rows.Close()
	}()

	// пробегаем по всем записям
	for rows.Next() {
		var farm Farm
		err = rows.Scan(&farm.ID, &farm.Name, &farm.Address)
		if err != nil {
			logger.Wr.Warn().Err(err).Msg("get farm instance error")
			return "", err
		}
		farms = append(farms, farm)
	}

	if err := rows.Err(); err != nil {
		logger.Wr.Warn().Err(err).Msg("get farm rows error")
		return "", err
	}

	if len(farms) == 0 {
		return "", NewEmptyError("no farm for current user")
	}
	data, err := json.Marshal(&farms)
	if err != nil {
		logger.Wr.Warn().Err(err).Msg("marshal to json error")
		return "", err
	}
	return string(data), nil
}

func (s *DBStorage) AddFarm(c context.Context, farm Farm) error {
	ctxIn, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	var userID int
	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1",
		FFarmID, TFarm, FAddress)
	row := s.db.QueryRowContext(ctxIn, sqlStr, farm.Address)
	err := row.Scan(&userID)

	if err == nil {
		logger.Wr.Info().Msg("farm already exist")
		return NewExistError("farm already exist")
	} else if err != sql.ErrNoRows {
		logger.Wr.Warn().Err(err).Msg("db request error")
		return err
	}

	//добавление
	sqlStr = fmt.Sprintf("INSERT INTO %s (%s, %s, %s) VALUES ($1, $2, $3)",
		TFarm, FName, FAddress, FUserID)
	_, err = s.db.ExecContext(ctxIn, sqlStr, farm.Name, farm.Address, farm.UserID)
	if err != nil {
		logger.Wr.Warn().Err(err).Msg("inserting farm error")
		return err
	}
	return nil
}

func (s *DBStorage) DelFarm(c context.Context, userID int, farmID int) error {
	ctxIn, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	sqlStr := fmt.Sprintf("UPDATE %s SET %s = TRUE "+
		" WHERE %s = $1 AND %s = $2 AND %s = FALSE",
		TFarm, FDeleted, FUserID, FFarmID, FDeleted)

	res, err := s.db.ExecContext(ctxIn, sqlStr, userID, farmID)
	if err != nil {
		logger.Wr.Warn().Err(err).Msg("db request error")
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		logger.Wr.Warn().Err(err).Msg("db request error")
		return err
	}
	if count == 0 {
		logger.Wr.Info().Msgf("no farm with index %s", farmID)
		return NewEmptyError("no farm for current user")
	}

	//TODO нужно обновлять таблицу коров, здоровья
	return nil
}

func (s *DBStorage) GetCowBreeds(c context.Context) (string, error) {
	ctxIn, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	var breeds []CowBreed
	sqlStr := fmt.Sprintf("SELECT %s, %s FROM %s",
		FBreedID, FBreed, TBreed)
	rows, err := s.db.QueryContext(ctxIn, sqlStr)

	if err != nil {
		logger.Wr.Warn().Err(err).Msg("db request error")
		return "", err
	}

	defer func() {
		_ = rows.Close()
	}()

	// пробегаем по всем записям
	for rows.Next() {
		var breed CowBreed
		err = rows.Scan(&breed.ID, &breed.Breed)
		if err != nil {
			logger.Wr.Warn().Err(err).Msg("get breed instance error")
			return "", err
		}
		breeds = append(breeds, breed)
	}

	if err := rows.Err(); err != nil {
		logger.Wr.Warn().Err(err).Msg("get breed rows error")
		return "", err
	}

	if len(breeds) == 0 {
		return "", NewEmptyError("no farm for current user")
	}
	data, err := json.Marshal(&breeds)
	if err != nil {
		logger.Wr.Warn().Err(err).Msg("marshal to json error")
		return "", err
	}
	return string(data), nil
}

//func (s *DBStorage) GetFarmInfo(ctx context.Context, farmID int) (string, error) {
//	return "", nil
//}
//
//func (s *DBStorage) GetCows(_ context.Context, farmID int) (string, error) {
//	return "", nil
//}

func (s *DBStorage) GetBolusesTypes(c context.Context) (string, error) {
	ctxIn, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var types []string
	sqlStr := "SELECT unnest(enum_range(NULL::bolus_type))"
	rows, err := s.db.QueryContext(ctxIn, sqlStr)

	if err != nil {
		logger.Wr.Warn().Err(err).Msg("db request error")
		return "", err
	}

	defer func() {
		_ = rows.Close()
	}()

	// пробегаем по всем записям
	for rows.Next() {
		var bolusType string
		err = rows.Scan(&bolusType)
		if err != nil {
			logger.Wr.Warn().Err(err).Msg("get bolus type error")
			return "", err
		}
		types = append(types, bolusType)
	}

	if err := rows.Err(); err != nil {
		logger.Wr.Warn().Err(err).Msg("get bolus types rows error")
		return "", err
	}

	if len(types) == 0 {
		return "", NewEmptyError("no bolus types")
	}
	data, err := json.Marshal(&types)
	if err != nil {
		logger.Wr.Warn().Err(err).Msg("marshal to json error")
		return "", err
	}
	return string(data), nil
}

//func (s *DBStorage) GetCowInfo(ctx context.Context, cowID int) (string, error) {
//	return "", nil
//}

//func (s *DBStorage) AddMonitoringData(ctx context.Context, data MonitoringData) error {
//	return nil
//}
//
//func (s *DBStorage) DeleteCows(ctx context.Context, IDs []int) error {
//	return nil
//}

func (s *DBStorage) AddCow(c context.Context, cow Cow) error {
	ctxIn, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	var bolusID int
	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1",
		FCowID, TCow, FBolus)
	row := s.db.QueryRowContext(ctxIn, sqlStr, cow.BolusNum)
	err := row.Scan(&bolusID)

	if err == nil {
		logger.Wr.Info().Msg("duplicate bolus")
		return NewExistError("duplicate bolus")
	}

	//добавление коровы
	sqlStr = fmt.Sprintf("INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING %s",
		TCow, FName, FBreedID, FFarmID, FBolus, FDateOfBorn, FAddedAt, FBolusType, FCowID)

	row = s.db.QueryRowContext(ctxIn, sqlStr, cow.Name, cow.BreedID, cow.FarmID,
		cow.BolusNum, cow.DateOfBorn, cow.AddedAt, cow.BolusType)

	err = row.Scan(&cow.ID)

	if err != nil {
		logger.Wr.Warn().Err(err).Msg("db request error")
		return err
	}

	//добавление в таблицу здоровье
	sqlStr = fmt.Sprintf("INSERT INTO %s(%s) VALUES ($1)",
		THealth, FCowID)
	_, err = s.db.ExecContext(ctxIn, sqlStr, cow.ID)
	if err != nil {
		logger.Wr.Warn().Err(err).Msg("creating health record error")
		return err
	}
	return nil
}

func (s *DBStorage) connect(DSN string) {
	var err error
	s.db, err = sql.Open("pgx", DSN)
	if err != nil {
		panic(err)
	}
}

func (s *DBStorage) createTables(ctx context.Context) {
	ctxIn, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// types
	sqlStr := fmt.Sprint("DO $$ BEGIN " +
		"IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'bolus_type') " +
		"THEN CREATE TYPE bolus_type AS ENUM ('С датчиком PH', 'Без датчика PH'); " +
		"END IF; END$$")
	_, err := s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msg("Fail then creating type bolus_type")
	}
	logger.Wr.Info().Msg("Type created: bolus_type")

	//1. user table
	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
		"(%s serial UNIQUE PRIMARY KEY, %s TEXT UNIQUE NOT NULL, %s TEXT NOT NULL)",
		TUser, FUserID, FLogin, FPassword)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
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
		"%s TEXT UNIQUE NOT NULL, %s INTEGER NOT NULL, "+
		"%s BOOLEAN NOT NULL DEFAULT FALSE)",
		TFarm, FFarmID, FName, FAddress, FUserID, FDeleted)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		logger.Wr.Panic().Err(err).Msgf("Fail then creating table %s", TFarm)
	}
	logger.Wr.Info().Msgf("Table created: %s", TFarm)

	//4. health table
	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
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
	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
		"(%s serial UNIQUE PRIMARY KEY, %s TEXT NOT NULL, "+
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
