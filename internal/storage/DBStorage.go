package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Sur0vy/cows_health.git/internal/helpers"
	"strconv"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/mattn/go-sqlite3"

	"github.com/Sur0vy/cows_health.git/internal/logger"
)

type DBStorage struct {
	log *logger.Logger
	db  *sql.DB
}

func NewDBStorage(ctx context.Context, DSN string, log *logger.Logger) *DBStorage {
	s := &DBStorage{
		log: log,
	}
	s.connect(DSN)
	s.createTables(ctx)
	return s
}

func (s *DBStorage) AddUser(c context.Context, user User) error {
	userHash := helpers.GetMD5Hash(user.Login)
	u := s.GetUser(c, userHash)
	if u != nil {
		s.log.Warn().Msgf("User %v already exists", user.Login)
		return NewExistError()
	}

	passwordHash, err := helpers.GetCryptoPassword(user.Password)
	if err != nil {
		s.log.Warn().Msg("Error than encrypting password")
		return err
	}

	ctxIn, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	//добавление
	sqlStr := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES ($1, $2)",
		TUser, FLogin, FPassword)
	_, err = s.db.ExecContext(ctxIn, sqlStr, userHash, passwordHash)
	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return err
	}
	return nil
}

func (s *DBStorage) GetUser(c context.Context, userHash string) *User {
	ctxIn, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	user := &User{}

	sqlStr := fmt.Sprintf("SELECT %s, %s FROM %s WHERE %s = $1",
		FUserID, FPassword, TUser, FLogin)
	row := s.db.QueryRowContext(ctxIn, sqlStr, userHash)
	err := row.Scan(&user.ID, &user.Password)

	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return nil
	}
	return user
}

func (s *DBStorage) GetUserHash(c context.Context, user User) (string, error) {
	userHash := helpers.GetMD5Hash(user.Login)
	u := s.GetUser(c, userHash)
	if u == nil {
		s.log.Warn().Msgf("User %v not exists", user.Login)
		return "", NewEmptyError()
	}

	if helpers.CheckPassword(u.Password, user.Password) {
		return userHash, nil
	}
	return "", NewEmptyError()
}

func (s *DBStorage) GetFarms(c context.Context, userID int) ([]Farm, error) {
	ctxIn, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var farms []Farm
	sqlStr := fmt.Sprintf("SELECT %s, %s, %s FROM %s WHERE %s = $1 AND NOT %s",
		FFarmID, FName, FAddress, TFarm, FUserID, FDeleted)
	rows, err := s.db.QueryContext(ctxIn, sqlStr, userID)

	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return farms, err
	}

	defer func() {
		_ = rows.Close()
	}()

	// пробегаем по всем записям
	for rows.Next() {
		var farm Farm
		err = rows.Scan(&farm.ID, &farm.Name, &farm.Address)
		if err != nil {
			s.log.Warn().Err(err).Msg("get farm instance error")
			return farms, err
		}
		farms = append(farms, farm)
	}

	if err := rows.Err(); err != nil {
		s.log.Warn().Err(err).Msg("get farm rows error")
		return farms, err
	}

	if len(farms) == 0 {
		return farms, NewEmptyError()
	}
	return farms, nil
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
		s.log.Info().Msg("farm already exist")
		return NewExistError()
	} else if err != sql.ErrNoRows {
		s.log.Warn().Err(err).Msg("db request error")
		return err
	}

	//добавление
	sqlStr = fmt.Sprintf("INSERT INTO %s (%s, %s, %s) VALUES ($1, $2, $3)",
		TFarm, FName, FAddress, FUserID)
	_, err = s.db.ExecContext(ctxIn, sqlStr, farm.Name, farm.Address, farm.UserID)
	if err != nil {
		s.log.Warn().Err(err).Msg("inserting farm error")
		return err
	}
	return nil
}

func (s *DBStorage) DelFarm(c context.Context, farmID int) error {
	ctxIn, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	sqlStr := fmt.Sprintf("UPDATE %s SET %s = TRUE "+
		" WHERE %s = $1 AND %s = FALSE",
		TFarm, FDeleted, FFarmID, FDeleted)

	res, err := s.db.ExecContext(ctxIn, sqlStr, farmID)
	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return err
	}
	if count == 0 {
		s.log.Info().Msgf("no farm with index %d", farmID)
		return NewEmptyError()
	}

	//TODO нужно обновлять таблицу коров, здоровья
	return nil
}

func (s *DBStorage) GetCowBreeds(c context.Context) (string, error) {
	ctxIn, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	var breeds []CowBreed
	sqlStr := fmt.Sprintf("SELECT %s, %s FROM %s",
		FBreedID, FName, TBreed)
	rows, err := s.db.QueryContext(ctxIn, sqlStr)

	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
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
			s.log.Warn().Err(err).Msg("get breed instance error")
			return "", err
		}
		breeds = append(breeds, breed)
	}

	if err := rows.Err(); err != nil {
		s.log.Warn().Err(err).Msg("get breed rows error")
		return "", err
	}

	if len(breeds) == 0 {
		return "", NewEmptyError()
	}
	data, err := json.Marshal(&breeds)
	if err != nil {
		s.log.Warn().Err(err).Msg("marshal to json error")
		return "", err
	}
	return string(data), nil
}

func (s *DBStorage) GetCows(c context.Context, farmID int) (string, error) {
	ctxIn, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	sqlStr := fmt.Sprintf("SELECT %s, %s, %s, %s, %s, %s, %s FROM %s "+
		"WHERE %s = $1 AND NOT %s",
		FCowID, FName, FBreedID, FBolus, FDateOfBorn, FAddedAt, FBolusType, TCow, FFarmID, FDeleted)
	rows, err := s.db.QueryContext(ctxIn, sqlStr, farmID)

	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return "", err
	}

	defer func() {
		_ = rows.Close()
	}()

	// пробегаем по всем записям
	var cows []Cow
	for rows.Next() {
		var cow Cow
		err = rows.Scan(&cow.ID, &cow.Name, &cow.BreedID, &cow.BolusNum,
			&cow.DateOfBorn, &cow.AddedAt, &cow.BolusType)
		if err != nil {
			s.log.Warn().Err(err).Msg("get cow instance error")
			return "", err
		}
		cow.FarmID = farmID
		cows = append(cows, cow)
	}

	if err := rows.Err(); err != nil {
		s.log.Warn().Err(err).Msg("get cow rows error")
		return "", err
	}

	if len(cows) == 0 {
		return "", NewEmptyError()
	}
	data, err := json.Marshal(&cows)
	if err != nil {
		s.log.Warn().Err(err).Msg("marshal to json error")
		return "", err
	}
	return string(data), nil
}

func (s *DBStorage) GetBolusesTypes(c context.Context) (string, error) {
	ctxIn, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var types []string
	sqlStr := "SELECT unnest(enum_range(NULL::bolus_type))"
	rows, err := s.db.QueryContext(ctxIn, sqlStr)

	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
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
			s.log.Warn().Err(err).Msg("get bolus type error")
			return "", err
		}
		types = append(types, bolusType)
	}

	if err := rows.Err(); err != nil {
		s.log.Warn().Err(err).Msg("get bolus types rows error")
		return "", err
	}

	if len(types) == 0 {
		return "", NewEmptyError()
	}
	data, err := json.Marshal(&types)
	if err != nil {
		s.log.Warn().Err(err).Msg("marshal to json error")
		return "", err
	}
	return string(data), nil
}

func (s *DBStorage) GetCowInfo(c context.Context, cowID int) (string, error) {
	ctxIn, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	var cowInfo CowInfo
	sqlStr := fmt.Sprintf("SELECT c.%s, b.%s, c.%s, c.%s, "+
		"c.%s, h.%s, h.%s, h.%s, "+
		"md.%s, md.%s, md.%s, md.%s, md.%s "+
		"FROM %s AS c "+
		"JOIN %s AS b ON b.%s = c.%s "+
		"JOIN %s AS h ON h.%s = c.%s "+
		"JOIN %s AS md ON md.%s = c.%s "+
		"WHERE ((EXTRACT(EPOCH FROM now()) - "+
		"EXTRACT(EPOCH FROM md.%s)) < $1) AND (md.%s = $2) "+
		"ORDER BY c.%s ASC",
		FName, FName, FBolus, FDateOfBorn,
		FBolusType, FEstrus, FIll, FUpdatedAt,
		FAddedAt, FPH, FTemperature, FMovement, FCharge,
		TCow,
		TBreed, FBreedID, FBreedID,
		THealth, FCowID, FCowID,
		TMonitoringData, FCowID, FCowID,
		FAddedAt, FCowID,
		FAddedAt)

	hour := 3600
	intervalInS := hour * 24
	rows, err := s.db.QueryContext(ctxIn, sqlStr, intervalInS, cowID)

	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return "", err
	}

	defer func() {
		_ = rows.Close()
	}()

	// пробегаем по всем записям
	for rows.Next() {
		var md MonitoringData
		err = rows.Scan(&cowInfo.Summary.Name, &cowInfo.Summary.Breed, &cowInfo.Summary.BolusNum,
			&cowInfo.Summary.DateOfBorn, &cowInfo.Summary.BolusType,
			&cowInfo.Health.Estrus, &cowInfo.Health.Ill, &cowInfo.Health.UpdatedAt,
			&md.AddedAt, &md.PH, &md.Temperature, &md.Movement, &md.Charge)
		if err != nil {
			s.log.Warn().Err(err).Msg("get monitoring data instance error")
			return "", err
		}
		cowInfo.History = append(cowInfo.History, md)
	}

	if err := rows.Err(); err != nil {
		s.log.Warn().Err(err).Msg("get cow info rows error")
		return "", err
	}

	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return "", err
	}

	data, err := json.Marshal(cowInfo)
	if err != nil {
		s.log.Warn().Err(err).Msg("marshal to json error")
		return "", err
	}
	return string(data), nil
}

func (s *DBStorage) HasBolus(c context.Context, BolusNum int) int {
	ctxIn, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	cowID := -1

	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1 LIMIT 1",
		FCowID, TCow, FBolus)
	row := s.db.QueryRowContext(ctxIn, sqlStr, BolusNum)
	err := row.Scan(&cowID)

	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return -1
	}
	return cowID
}

func (s *DBStorage) AddMonitoringData(c context.Context, data MonitoringData) error {
	ctxIn, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	//добавление коровы
	sqlStr := fmt.Sprintf("INSERT INTO %s (%s, %s, %s, %s, %s, %s) "+
		"VALUES ($1, $2, $3, $4, $5, $6)",
		TMonitoringData, FCowID, FAddedAt, FPH, FTemperature, FMovement, FCharge)

	_, err := s.db.ExecContext(ctxIn, sqlStr, data.CowID, data.AddedAt, data.PH,
		data.Temperature, data.Movement, data.Charge)
	if err != nil {
		s.log.Warn().Err(err).Msg("inserting monitoring data error")
		return err
	}
	return nil
}

func (s *DBStorage) GetMonitoringData(c context.Context, cowID int, interval int) ([]MonitoringData, error) {
	ctxIn, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	var res []MonitoringData
	sqlStr := fmt.Sprintf("SELECT %s, %s, %s, %s FROM %s "+
		"WHERE ((EXTRACT(EPOCH FROM now()) - "+
		"EXTRACT(EPOCH FROM %s)) < $1) AND (%s = $2)",
		FTemperature, FMovement, FPH, FAddedAt, TMonitoringData, FAddedAt, FCowID)

	//now := time.Now()
	min := 60
	intervalInS := min * interval
	rows, err := s.db.QueryContext(ctxIn, sqlStr, intervalInS, cowID)

	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return res, err
	}

	defer func() {
		_ = rows.Close()
	}()

	// пробегаем по всем записям
	for rows.Next() {
		var md MonitoringData
		err = rows.Scan(&md.Temperature, &md.Movement, &md.PH, &md.AddedAt)
		if err != nil {
			s.log.Warn().Err(err).Msg("get monitoring data instance error")
			return nil, err
		}
		res = append(res, md)
	}

	if err := rows.Err(); err != nil {
		s.log.Warn().Err(err).Msg("get monitoring data rows error")
		return nil, err
	}

	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return nil, err
	}

	return res, nil
}

func (s *DBStorage) UpdateHealth(c context.Context, data Health) error {
	ctxIn, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	sqlStr := fmt.Sprintf("UPDATE %s SET %s = $1, %s = $2, %s = $3 WHERE %s = $4",
		THealth, FUpdatedAt, FIll, FEstrus, FCowID)

	_, err := s.db.ExecContext(ctxIn, sqlStr, data.UpdatedAt, data.Ill, data.Estrus, data.CowID)
	if err != nil {
		s.log.Warn().Err(err).Msg("inserting health data error")
		return err
	}
	return nil
}

func (s *DBStorage) AddCow(c context.Context, cow Cow) error {
	ctxIn, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	var bolusID int
	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1",
		FCowID, TCow, FBolus)
	row := s.db.QueryRowContext(ctxIn, sqlStr, cow.BolusNum)
	err := row.Scan(&bolusID)

	if err == nil {
		s.log.Info().Msg("duplicate bolus")
		return NewExistError()
	}

	//добавление коровы
	sqlStr = fmt.Sprintf("INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING %s",
		TCow, FName, FBreedID, FFarmID, FBolus, FDateOfBorn, FAddedAt, FBolusType, FCowID)

	row = s.db.QueryRowContext(ctxIn, sqlStr, cow.Name, cow.BreedID, cow.FarmID,
		cow.BolusNum, cow.DateOfBorn, cow.AddedAt, cow.BolusType)

	err = row.Scan(&cow.ID)

	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return err
	}

	//добавление в таблицу здоровье
	sqlStr = fmt.Sprintf("INSERT INTO %s(%s) VALUES ($1)",
		THealth, FCowID)
	_, err = s.db.ExecContext(ctxIn, sqlStr, cow.ID)
	if err != nil {
		s.log.Warn().Err(err).Msg("creating health record error")
		return err
	}
	return nil
}

func (s *DBStorage) DeleteCows(c context.Context, CowIDs []int) error {
	ctxIn, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var arr []interface{}
	var pos strings.Builder

	for num, ID := range CowIDs {
		if pos.Len() != 0 {
			pos.WriteString(", ")
		}
		arr = append(arr, ID)
		pos.WriteString("$")
		pos.WriteString(strconv.Itoa(num + 1))
	}

	sqlStr := fmt.Sprintf("UPDATE %s SET %s = TRUE "+
		" WHERE %s IN("+pos.String()+") AND %s = FALSE",
		TCow, FDeleted, FCowID, FDeleted)

	res, err := s.db.ExecContext(ctxIn, sqlStr, arr...)
	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return err
	}
	if count == 0 {
		s.log.Info().Msgf("no cows with indexes %v", CowIDs)
		return NewEmptyError()
	}

	//health
	sqlStr = fmt.Sprintf("UPDATE %s SET %s = TRUE "+
		" WHERE %s IN("+pos.String()+") AND %s = FALSE",
		THealth, FDeleted, FCowID, FDeleted)

	res, err = s.db.ExecContext(ctxIn, sqlStr, arr...)
	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return err
	}
	count, err = res.RowsAffected()
	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return err
	}
	if count == 0 {
		s.log.Info().Msgf("no health with indexes %v", CowIDs)
		return NewEmptyError()
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
		s.log.Panic().Err(err).Msg("Fail then creating type bolus_type")
	}
	s.log.Info().Msg("Type created: bolus_type")

	//1. user table
	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
		"(%s serial UNIQUE PRIMARY KEY, %s TEXT UNIQUE NOT NULL, %s TEXT NOT NULL)",
		TUser, FUserID, FLogin, FPassword)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		s.log.Panic().Err(err).Msgf("Fail then creating table %s", TUser)
	}
	s.log.Info().Msgf("Table created: %s", TUser)

	//2. breed table
	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
		"(%s serial UNIQUE PRIMARY KEY, %s TEXT NOT NULL)",
		TBreed, FBreedID, FName)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		s.log.Panic().Err(err).Msgf("Fail then creating table %s", TBreed)
	}
	s.log.Info().Msgf("Table created: %s", TBreed)

	//3. farm table
	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
		"(%s serial UNIQUE PRIMARY KEY, %s TEXT NOT NULL, "+
		"%s TEXT UNIQUE NOT NULL, %s INTEGER NOT NULL, "+
		"%s BOOLEAN NOT NULL DEFAULT FALSE)",
		TFarm, FFarmID, FName, FAddress, FUserID, FDeleted)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		s.log.Panic().Err(err).Msgf("Fail then creating table %s", TFarm)
	}
	s.log.Info().Msgf("Table created: %s", TFarm)

	//4. health table
	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
		"(%s INTEGER UNIQUE PRIMARY KEY, %s BOOLEAN, "+
		"%s TEXT, %s TIMESTAMP with time zone, %s BOOLEAN NOT NULL DEFAULT FALSE)",
		THealth, FCowID, FEstrus, FIll, FUpdatedAt, FDeleted)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		s.log.Panic().Err(err).Msgf("Fail then creating table %s", THealth)
	}
	s.log.Info().Msgf("Table created: %s", THealth)

	//5. monitoring data table
	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
		"(%s serial UNIQUE PRIMARY KEY, %s INTEGER NOT NULL, "+
		"%s TIMESTAMP with time zone, %s FLOAT, %s FLOAT, %s FLOAT, %s FLOAT)",
		TMonitoringData, FMDID, FCowID, FAddedAt, FPH, FTemperature, FMovement, FCharge)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		s.log.Panic().Err(err).Msgf("Fail then creating table %s", TMonitoringData)
	}
	s.log.Info().Msgf("Table created: %s", TMonitoringData)

	//6. cow table (проверить, есть ли тип такой)
	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
		"(%s serial UNIQUE PRIMARY KEY, %s TEXT NOT NULL, "+
		"%s INTEGER NOT NULL, %s INTEGER NOT NULL, "+
		"%s INTEGER UNIQUE NOT NULL, %s DATE NOT NULL, "+
		"%s TIMESTAMP with time zone NOT NULL, %s bolus_type NOT NULL, "+
		"%s BOOLEAN NOT NULL DEFAULT FALSE)",
		TCow, FCowID, FName, FBreedID, FFarmID, FBolus, FDateOfBorn, FAddedAt, FBolusType, FDeleted)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		s.log.Panic().Err(err).Msgf("Fail then creating table %s", TCow)
	}
	s.log.Info().Msgf("Table created: %s", TCow)

	//links
	//user-farm link
	sqlStr = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_user_farm; "+
		"ALTER TABLE %s ADD CONSTRAINT fk_user_farm FOREIGN KEY (%s) REFERENCES %s (%s)",
		TFarm, TFarm, FUserID, TUser, FUserID)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		s.log.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", TFarm, TUser)
	}
	s.log.Info().Msgf("foreign key created: %s <-> %s", TFarm, TUser)

	//cow-breed link
	sqlStr = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_cow_breed; "+
		"ALTER TABLE %s ADD CONSTRAINT fk_cow_breed FOREIGN KEY (%s) REFERENCES %s (%s)",
		TCow, TCow, FBreedID, TBreed, FBreedID)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		s.log.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", TCow, TBreed)
	}
	s.log.Info().Msgf("foreign key created: %s <-> %s", TCow, TBreed)

	//cow-farm link
	sqlStr = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_cow_farm; "+
		"ALTER TABLE %s ADD CONSTRAINT fk_cow_farm FOREIGN KEY (%s) REFERENCES %s (%s)",
		TCow, TCow, FFarmID, TFarm, FFarmID)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		s.log.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", TCow, TFarm)
	}
	s.log.Info().Msgf("foreign key created: %s <-> %s", TCow, TFarm)

	//health-cow link
	sqlStr = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_cow_health; "+
		"ALTER TABLE %s ADD CONSTRAINT fk_cow_health FOREIGN KEY (%s) REFERENCES %s (%s)",
		THealth, THealth, FCowID, TCow, FCowID)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		s.log.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", THealth, TCow)
	}
	s.log.Info().Msgf("foreign key created: %s <-> %s", THealth, TCow)

	//monitoring data-cow link
	sqlStr = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_cow_md; "+
		"ALTER TABLE %s ADD CONSTRAINT fk_cow_md FOREIGN KEY (%s) REFERENCES %s (%s)",
		TMonitoringData, TMonitoringData, FCowID, TCow, FCowID)
	_, err = s.db.ExecContext(ctxIn, sqlStr)
	if err != nil {
		s.log.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", TMonitoringData, TCow)
	}
	s.log.Info().Msgf("foreign key created: %s <-> %s", TMonitoringData, TCow)
}

func (s *DBStorage) Ping() error {
	return s.db.Ping()
}
