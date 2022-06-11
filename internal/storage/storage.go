package storage

import (
	"context"
	"crypto/md5"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	AddUser(ctx context.Context, user User) (string, error)
	GetUserHash(ctx context.Context, user User) (string, error)
	GetUser(ctx context.Context, userHash string) *User
	//GetUser(ctx context.Context, cookie string) (int, error)
	//
	//GetFarms(ctx context.Context, userID int) (string, error)
	//AddFarm(ctx context.Context, farm Farm) error
	//GetFarmInfo(ctx context.Context, farmID int) (string, error)
	//
	//GetCows(ctx context.Context, farmID int) (string, error)
	//GetCowInfo(ctx context.Context, farmID int) (string, error)
	//GetCowBreeds(ctx context.Context) (string, error)
	//DeleteCows(ctx context.Context, IDs []int) error
	//AddCow(ctx context.Context, cow Cow) error
	//
	//GetBolusesTypes(ctx context.Context) (string, error)
	//AddMonitoringData(ctx context.Context, data MonitoringData) error
}

func getMD5Hash(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

func getCryptoPassword(text string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hashedPassword), nil
}

func checkPassword(hash string, password string) bool {
	h, err := hex.DecodeString(hash)
	if err == nil {
		err = bcrypt.CompareHashAndPassword(h, []byte(password))
		if err != nil {
			return false
		}
	}
	return true
}
