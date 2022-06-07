package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type DBStorage struct {
	//userStorage users.UserStorage
	//db          *sql.DB
}

func NewDBStorage(_ context.Context, DSN string) Storage {
	s := &DBStorage{}
	//s.Connect()
	//s.CreateTables(ctx)
	//s.userStorage = users.NewDBUserStorage(s.db)
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
