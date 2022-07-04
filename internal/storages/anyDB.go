package storages

//func (s *DBStorage) GetCowBreeds(c context.Context) ([]models.CowBreed, error) {
//	ctxIn, cancel := context.WithTimeout(c, 2*time.Second)
//	defer cancel()
//
//	var breeds []models.CowBreed
//	sqlStr := fmt.Sprintf("SELECT %s, %s FROM %s",
//		internal.FBreedID, internal.FName, internal.TBreed)
//	rows, err := s.db.QueryContext(ctxIn, sqlStr)
//
//	if err != nil {
//		s.log.Warn().Err(err).Msg("db request error")
//		return breeds, err
//	}
//
//	defer func() {
//		_ = rows.Close()
//	}()
//
//	// пробегаем по всем записям
//	for rows.Next() {
//		var breed models.CowBreed
//		err = rows.Scan(&breed.ID, &breed.Breed)
//		if err != nil {
//			s.log.Warn().Err(err).Msg("get breed instance error")
//			return breeds, err
//		}
//		breeds = append(breeds, breed)
//	}
//
//	if err := rows.Err(); err != nil {
//		s.log.Warn().Err(err).Msg("get breed rows error")
//		return breeds, err
//	}
//
//	return breeds, nil
//}
//
//func (s *DBStorage) GetCows(c context.Context, farmID int) (string, error) {
//	ctxIn, cancel := context.WithTimeout(c, 5*time.Second)
//	defer cancel()
//
//	sqlStr := fmt.Sprintf("SELECT %s, %s, %s, %s, %s, %s, %s FROM %s "+
//		"WHERE %s = $1 AND NOT %s",
//		internal.FCowID, internal.FName, internal.FBreedID, internal.FBolus, internal.FDateOfBorn, internal.FAddedAt, internal.FBolusType, internal.TCow, internal.FFarmID, entity.FDeleted)
//	rows, err := s.db.QueryContext(ctxIn, sqlStr, farmID)
//
//	if err != nil {
//		s.log.Warn().Err(err).Msg("db request error")
//		return "", err
//	}
//
//	defer func() {
//		_ = rows.Close()
//	}()
//
//	// пробегаем по всем записям
//	var cows []models.Cow
//	for rows.Next() {
//		var cow models.Cow
//		err = rows.Scan(&cow.ID, &cow.Name, &cow.BreedID, &cow.BolusNum,
//			&cow.DateOfBorn, &cow.AddedAt, &cow.BolusType)
//		if err != nil {
//			s.log.Warn().Err(err).Msg("get cow instance error")
//			return "", err
//		}
//		cow.FarmID = farmID
//		cows = append(cows, cow)
//	}
//
//	if err := rows.Err(); err != nil {
//		s.log.Warn().Err(err).Msg("get cow rows error")
//		return "", err
//	}
//
//	if len(cows) == 0 {
//		return "", errors.NewEmptyError()
//	}
//	data, err := json.Marshal(&cows)
//	if err != nil {
//		s.log.Warn().Err(err).Msg("marshal to json error")
//		return "", err
//	}
//	return string(data), nil
//}
//
//func (s *DBStorage) GetBolusesTypes(c context.Context) (string, error) {
//	ctxIn, cancel := context.WithTimeout(c, 5*time.Second)
//	defer cancel()
//
//	var types []string
//	sqlStr := "SELECT unnest(enum_range(NULL::bolus_type))"
//	rows, err := s.db.QueryContext(ctxIn, sqlStr)
//
//	if err != nil {
//		s.log.Warn().Err(err).Msg("db request error")
//		return "", err
//	}
//
//	defer func() {
//		_ = rows.Close()
//	}()
//
//	// пробегаем по всем записям
//	for rows.Next() {
//		var bolusType string
//		err = rows.Scan(&bolusType)
//		if err != nil {
//			s.log.Warn().Err(err).Msg("get bolus type error")
//			return "", err
//		}
//		types = append(types, bolusType)
//	}
//
//	if err := rows.Err(); err != nil {
//		s.log.Warn().Err(err).Msg("get bolus types rows error")
//		return "", err
//	}
//
//	if len(types) == 0 {
//		return "", errors.NewEmptyError()
//	}
//	data, err := json.Marshal(&types)
//	if err != nil {
//		s.log.Warn().Err(err).Msg("marshal to json error")
//		return "", err
//	}
//	return string(data), nil
//}
//
//func (s *DBStorage) GetCowInfo(c context.Context, cowID int) (string, error) {
//	ctxIn, cancel := context.WithTimeout(c, time.Second)
//	defer cancel()
//
//	var cowInfo models.CowInfo
//	sqlStr := fmt.Sprintf("SELECT c.%s, b.%s, c.%s, c.%s, "+
//		"c.%s, h.%s, h.%s, h.%s, "+
//		"md.%s, md.%s, md.%s, md.%s, md.%s "+
//		"FROM %s AS c "+
//		"JOIN %s AS b ON b.%s = c.%s "+
//		"JOIN %s AS h ON h.%s = c.%s "+
//		"JOIN %s AS md ON md.%s = c.%s "+
//		"WHERE ((EXTRACT(EPOCH FROM now()) - "+
//		"EXTRACT(EPOCH FROM md.%s)) < $1) AND (md.%s = $2) "+
//		"ORDER BY c.%s ASC",
//		internal.FName, internal.FName, internal.FBolus, internal.FDateOfBorn,
//		internal.FBolusType, internal.FEstrus, internal.FIll, internal.FUpdatedAt,
//		internal.FAddedAt, internal.FPH, internal.FTemperature, internal.FMovement, internal.FCharge,
//		internal.TCow,
//		internal.TBreed, internal.FBreedID, internal.FBreedID,
//		internal.THealth, internal.FCowID, internal.FCowID,
//		internal.TMonitoringData, internal.FCowID, internal.FCowID,
//		internal.FAddedAt, internal.FCowID,
//		internal.FAddedAt)
//
//	hour := 3600
//	intervalInS := hour * 24
//	rows, err := s.db.QueryContext(ctxIn, sqlStr, intervalInS, cowID)
//
//	if err != nil {
//		s.log.Warn().Err(err).Msg("db request error")
//		return "", err
//	}
//
//	defer func() {
//		_ = rows.Close()
//	}()
//
//	// пробегаем по всем записям
//	for rows.Next() {
//		var md models.MonitoringData
//		err = rows.Scan(&cowInfo.Summary.Name, &cowInfo.Summary.Breed, &cowInfo.Summary.BolusNum,
//			&cowInfo.Summary.DateOfBorn, &cowInfo.Summary.BolusType,
//			&cowInfo.Health.Estrus, &cowInfo.Health.Ill, &cowInfo.Health.UpdatedAt,
//			&md.AddedAt, &md.PH, &md.Temperature, &md.Movement, &md.Charge)
//		if err != nil {
//			s.log.Warn().Err(err).Msg("get monitoring data instance error")
//			return "", err
//		}
//		cowInfo.History = append(cowInfo.History, md)
//	}
//
//	if err := rows.Err(); err != nil {
//		s.log.Warn().Err(err).Msg("get cow info rows error")
//		return "", err
//	}
//
//	if err != nil {
//		s.log.Warn().Err(err).Msg("db request error")
//		return "", err
//	}
//
//	data, err := json.Marshal(cowInfo)
//	if err != nil {
//		s.log.Warn().Err(err).Msg("marshal to json error")
//		return "", err
//	}
//	return string(data), nil
//}
//
//func (s *DBStorage) HasBolus(c context.Context, BolusNum int) int {
//	ctxIn, cancel := context.WithTimeout(c, time.Second)
//	defer cancel()
//
//	cowID := -1
//
//	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1 LIMIT 1",
//		internal.FCowID, internal.TCow, internal.FBolus)
//	row := s.db.QueryRowContext(ctxIn, sqlStr, BolusNum)
//	err := row.Scan(&cowID)
//
//	if err != nil {
//		s.log.Warn().Err(err).Msg("db request error")
//		return -1
//	}
//	return cowID
//}
//
//func (s *DBStorage) AddMonitoringData(c context.Context, data models.MonitoringData) error {
//	ctxIn, cancel := context.WithTimeout(c, time.Second)
//	defer cancel()
//
//	//добавление коровы
//	sqlStr := fmt.Sprintf("INSERT INTO %s (%s, %s, %s, %s, %s, %s) "+
//		"VALUES ($1, $2, $3, $4, $5, $6)",
//		internal.TMonitoringData, internal.FCowID, internal.FAddedAt, internal.FPH, internal.FTemperature, internal.FMovement, internal.FCharge)
//
//	_, err := s.db.ExecContext(ctxIn, sqlStr, data.CowID, data.AddedAt, data.PH,
//		data.Temperature, data.Movement, data.Charge)
//	if err != nil {
//		s.log.Warn().Err(err).Msg("inserting monitoring data error")
//		return err
//	}
//	return nil
//}
//
//func (s *DBStorage) GetMonitoringData(c context.Context, cowID int, interval int) ([]models.MonitoringData, error) {
//	ctxIn, cancel := context.WithTimeout(c, time.Second)
//	defer cancel()
//
//	var res []models.MonitoringData
//	sqlStr := fmt.Sprintf("SELECT %s, %s, %s, %s FROM %s "+
//		"WHERE ((EXTRACT(EPOCH FROM now()) - "+
//		"EXTRACT(EPOCH FROM %s)) < $1) AND (%s = $2)",
//		internal.FTemperature, internal.FMovement, internal.FPH, internal.FAddedAt, internal.TMonitoringData, internal.FAddedAt, internal.FCowID)
//
//	//now := time.Now()
//	min := 60
//	intervalInS := min * interval
//	rows, err := s.db.QueryContext(ctxIn, sqlStr, intervalInS, cowID)
//
//	if err != nil {
//		s.log.Warn().Err(err).Msg("db request error")
//		return res, err
//	}
//
//	defer func() {
//		_ = rows.Close()
//	}()
//
//	// пробегаем по всем записям
//	for rows.Next() {
//		var md models.MonitoringData
//		err = rows.Scan(&md.Temperature, &md.Movement, &md.PH, &md.AddedAt)
//		if err != nil {
//			s.log.Warn().Err(err).Msg("get monitoring data instance error")
//			return nil, err
//		}
//		res = append(res, md)
//	}
//
//	if err := rows.Err(); err != nil {
//		s.log.Warn().Err(err).Msg("get monitoring data rows error")
//		return nil, err
//	}
//
//	if err != nil {
//		s.log.Warn().Err(err).Msg("db request error")
//		return nil, err
//	}
//
//	return res, nil
//}
//
//func (s *DBStorage) UpdateHealth(c context.Context, data models.Health) error {
//	ctxIn, cancel := context.WithTimeout(c, time.Second)
//	defer cancel()
//
//	sqlStr := fmt.Sprintf("UPDATE %s SET %s = $1, %s = $2, %s = $3 WHERE %s = $4",
//		internal.THealth, internal.FUpdatedAt, internal.FIll, internal.FEstrus, internal.FCowID)
//
//	_, err := s.db.ExecContext(ctxIn, sqlStr, data.UpdatedAt, data.Ill, data.Estrus, data.CowID)
//	if err != nil {
//		s.log.Warn().Err(err).Msg("inserting health data error")
//		return err
//	}
//	return nil
//}
//
//func (s *DBStorage) AddCow(c context.Context, cow models.Cow) error {
//	ctxIn, cancel := context.WithTimeout(c, 2*time.Second)
//	defer cancel()
//
//	var bolusID int
//	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1",
//		internal.FCowID, internal.TCow, internal.FBolus)
//	row := s.db.QueryRowContext(ctxIn, sqlStr, cow.BolusNum)
//	err := row.Scan(&bolusID)
//
//	if err == nil {
//		s.log.Info().Msg("duplicate bolus")
//		return errors.NewExistError()
//	}
//
//	//добавление коровы
//	sqlStr = fmt.Sprintf("INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s) "+
//		"VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING %s",
//		internal.TCow, internal.FName, internal.FBreedID, internal.FFarmID, internal.FBolus, internal.FDateOfBorn, internal.FAddedAt, internal.FBolusType, internal.FCowID)
//
//	row = s.db.QueryRowContext(ctxIn, sqlStr, cow.Name, cow.BreedID, cow.FarmID,
//		cow.BolusNum, cow.DateOfBorn, cow.AddedAt, cow.BolusType)
//
//	err = row.Scan(&cow.ID)
//
//	if err != nil {
//		s.log.Warn().Err(err).Msg("db request error")
//		return err
//	}
//
//	//добавление в таблицу здоровье
//	sqlStr = fmt.Sprintf("INSERT INTO %s(%s) VALUES ($1)",
//		internal.THealth, internal.FCowID)
//	_, err = s.db.ExecContext(ctxIn, sqlStr, cow.ID)
//	if err != nil {
//		s.log.Warn().Err(err).Msg("creating health record error")
//		return err
//	}
//	return nil
//}
//
//func (s *DBStorage) DeleteCows(c context.Context, CowIDs []int) error {
//	ctxIn, cancel := context.WithTimeout(c, 5*time.Second)
//	defer cancel()
//
//	var arr []interface{}
//	var pos strings.Builder
//
//	for num, ID := range CowIDs {
//		if pos.Len() != 0 {
//			pos.WriteString(", ")
//		}
//		arr = append(arr, ID)
//		pos.WriteString("$")
//		pos.WriteString(strconv.Itoa(num + 1))
//	}
//
//	sqlStr := fmt.Sprintf("UPDATE %s SET %s = TRUE "+
//		" WHERE %s IN("+pos.String()+") AND %s = FALSE",
//		internal.TCow, entity.FDeleted, internal.FCowID, entity.FDeleted)
//
//	res, err := s.db.ExecContext(ctxIn, sqlStr, arr...)
//	if err != nil {
//		s.log.Warn().Err(err).Msg("db request error")
//		return err
//	}
//	count, err := res.RowsAffected()
//	if err != nil {
//		s.log.Warn().Err(err).Msg("db request error")
//		return err
//	}
//	if count == 0 {
//		s.log.Info().Msgf("no cows with indexes %v", CowIDs)
//		return errors.NewEmptyError()
//	}
//
//	//health
//	sqlStr = fmt.Sprintf("UPDATE %s SET %s = TRUE "+
//		" WHERE %s IN("+pos.String()+") AND %s = FALSE",
//		internal.THealth, entity.FDeleted, internal.FCowID, entity.FDeleted)
//
//	res, err = s.db.ExecContext(ctxIn, sqlStr, arr...)
//	if err != nil {
//		s.log.Warn().Err(err).Msg("db request error")
//		return err
//	}
//	count, err = res.RowsAffected()
//	if err != nil {
//		s.log.Warn().Err(err).Msg("db request error")
//		return err
//	}
//	if count == 0 {
//		s.log.Info().Msgf("no health with indexes %v", CowIDs)
//		return errors.NewEmptyError()
//	}
//	return nil
//}
//
//func (s *DBStorage) connect(DSN string) {
//	var err error
//	s.db, err = sql.Open("pgx", DSN)
//	if err != nil {
//		panic(err)
//	}
//}
//
//func (s *DBStorage) Ping() error {
//	return s.db.Ping()
//}
//
//func (s *DBStorage) createTables(ctx context.Context) {
//	ctxIn, cancel := context.WithTimeout(ctx, 5*time.Second)
//	defer cancel()
//
//	// types
//	sqlStr := fmt.Sprint("DO $$ BEGIN " +
//		"IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'bolus_type') " +
//		"THEN CREATE TYPE bolus_type AS ENUM ('С датчиком PH', 'Без датчика PH'); " +
//		"END IF; END$$")
//	_, err := s.db.ExecContext(ctxIn, sqlStr)
//	if err != nil {
//		s.log.Panic().Err(err).Msg("Fail then creating type bolus_type")
//	}
//	s.log.Info().Msg("Type created: bolus_type")
//
//	//1. user table
//	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
//		"(%s serial UNIQUE PRIMARY KEY, %s TEXT UNIQUE NOT NULL, %s TEXT NOT NULL)",
//		entity.TUser, entity.FUserID, entity.FLogin, entity.FPassword)
//	_, err = s.db.ExecContext(ctxIn, sqlStr)
//	if err != nil {
//		s.log.Panic().Err(err).Msgf("Fail then creating table %s", entity.TUser)
//	}
//	s.log.Info().Msgf("Table created: %s", entity.TUser)
//
//	//2. breed table
//	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
//		"(%s serial UNIQUE PRIMARY KEY, %s TEXT NOT NULL)",
//		internal.TBreed, internal.FBreedID, internal.FName)
//	_, err = s.db.ExecContext(ctxIn, sqlStr)
//	if err != nil {
//		s.log.Panic().Err(err).Msgf("Fail then creating table %s", internal.TBreed)
//	}
//	s.log.Info().Msgf("Table created: %s", internal.TBreed)
//
//	//3. farm table
//	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
//		"(%s serial UNIQUE PRIMARY KEY, %s TEXT NOT NULL, "+
//		"%s TEXT UNIQUE NOT NULL, %s INTEGER NOT NULL, "+
//		"%s BOOLEAN NOT NULL DEFAULT FALSE)",
//		internal.TFarm, internal.FFarmID, internal.FName, internal.FAddress, entity.FUserID, entity.FDeleted)
//	_, err = s.db.ExecContext(ctxIn, sqlStr)
//	if err != nil {
//		s.log.Panic().Err(err).Msgf("Fail then creating table %s", internal.TFarm)
//	}
//	s.log.Info().Msgf("Table created: %s", internal.TFarm)
//
//	//4. health table
//	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
//		"(%s INTEGER UNIQUE PRIMARY KEY, %s BOOLEAN, "+
//		"%s TEXT, %s TIMESTAMP with time zone, %s BOOLEAN NOT NULL DEFAULT FALSE)",
//		internal.THealth, internal.FCowID, internal.FEstrus, internal.FIll, internal.FUpdatedAt, entity.FDeleted)
//	_, err = s.db.ExecContext(ctxIn, sqlStr)
//	if err != nil {
//		s.log.Panic().Err(err).Msgf("Fail then creating table %s", internal.THealth)
//	}
//	s.log.Info().Msgf("Table created: %s", internal.THealth)
//
//	//5. monitoring data table
//	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
//		"(%s serial UNIQUE PRIMARY KEY, %s INTEGER NOT NULL, "+
//		"%s TIMESTAMP with time zone, %s FLOAT, %s FLOAT, %s FLOAT, %s FLOAT)",
//		internal.TMonitoringData, internal.FMDID, internal.FCowID, internal.FAddedAt, internal.FPH, internal.FTemperature, internal.FMovement, internal.FCharge)
//	_, err = s.db.ExecContext(ctxIn, sqlStr)
//	if err != nil {
//		s.log.Panic().Err(err).Msgf("Fail then creating table %s", internal.TMonitoringData)
//	}
//	s.log.Info().Msgf("Table created: %s", internal.TMonitoringData)
//
//	//6. cow table (проверить, есть ли тип такой)
//	sqlStr = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
//		"(%s serial UNIQUE PRIMARY KEY, %s TEXT NOT NULL, "+
//		"%s INTEGER NOT NULL, %s INTEGER NOT NULL, "+
//		"%s INTEGER UNIQUE NOT NULL, %s DATE NOT NULL, "+
//		"%s TIMESTAMP with time zone NOT NULL, %s bolus_type NOT NULL, "+
//		"%s BOOLEAN NOT NULL DEFAULT FALSE)",
//		internal.TCow, internal.FCowID, internal.FName, internal.FBreedID, internal.FFarmID, internal.FBolus, internal.FDateOfBorn, internal.FAddedAt, internal.FBolusType, entity.FDeleted)
//	_, err = s.db.ExecContext(ctxIn, sqlStr)
//	if err != nil {
//		s.log.Panic().Err(err).Msgf("Fail then creating table %s", internal.TCow)
//	}
//	s.log.Info().Msgf("Table created: %s", internal.TCow)
//
//	//links
//	//user-farm link
//	sqlStr = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_user_farm; "+
//		"ALTER TABLE %s ADD CONSTRAINT fk_user_farm FOREIGN KEY (%s) REFERENCES %s (%s)",
//		internal.TFarm, internal.TFarm, entity.FUserID, entity.TUser, entity.FUserID)
//	_, err = s.db.ExecContext(ctxIn, sqlStr)
//	if err != nil {
//		s.log.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", internal.TFarm, entity.TUser)
//	}
//	s.log.Info().Msgf("foreign key created: %s <-> %s", internal.TFarm, entity.TUser)
//
//	//cow-breed link
//	sqlStr = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_cow_breed; "+
//		"ALTER TABLE %s ADD CONSTRAINT fk_cow_breed FOREIGN KEY (%s) REFERENCES %s (%s)",
//		internal.TCow, internal.TCow, internal.FBreedID, internal.TBreed, internal.FBreedID)
//	_, err = s.db.ExecContext(ctxIn, sqlStr)
//	if err != nil {
//		s.log.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", internal.TCow, internal.TBreed)
//	}
//	s.log.Info().Msgf("foreign key created: %s <-> %s", internal.TCow, internal.TBreed)
//
//	//cow-farm link
//	sqlStr = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_cow_farm; "+
//		"ALTER TABLE %s ADD CONSTRAINT fk_cow_farm FOREIGN KEY (%s) REFERENCES %s (%s)",
//		internal.TCow, internal.TCow, internal.FFarmID, internal.TFarm, internal.FFarmID)
//	_, err = s.db.ExecContext(ctxIn, sqlStr)
//	if err != nil {
//		s.log.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", internal.TCow, internal.TFarm)
//	}
//	s.log.Info().Msgf("foreign key created: %s <-> %s", internal.TCow, internal.TFarm)
//
//	//health-cow link
//	sqlStr = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_cow_health; "+
//		"ALTER TABLE %s ADD CONSTRAINT fk_cow_health FOREIGN KEY (%s) REFERENCES %s (%s)",
//		internal.THealth, internal.THealth, internal.FCowID, internal.TCow, internal.FCowID)
//	_, err = s.db.ExecContext(ctxIn, sqlStr)
//	if err != nil {
//		s.log.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", internal.THealth, internal.TCow)
//	}
//	s.log.Info().Msgf("foreign key created: %s <-> %s", internal.THealth, internal.TCow)
//
//	//monitoring data-cow link
//	sqlStr = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS fk_cow_md; "+
//		"ALTER TABLE %s ADD CONSTRAINT fk_cow_md FOREIGN KEY (%s) REFERENCES %s (%s)",
//		internal.TMonitoringData, internal.TMonitoringData, internal.FCowID, internal.TCow, internal.FCowID)
//	_, err = s.db.ExecContext(ctxIn, sqlStr)
//	if err != nil {
//		s.log.Panic().Err(err).Msgf("Fail then creating foreign key  %s <-> %s", internal.TMonitoringData, internal.TCow)
//	}
//	s.log.Info().Msgf("foreign key created: %s <-> %s", internal.TMonitoringData, internal.TCow)
//}
