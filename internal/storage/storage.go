package storage

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	AddUser(с context.Context, user User) (string, error)
	GetUserHash(с context.Context, user User) (string, error)
	GetUser(с context.Context, userHash string) *User

	GetFarms(с context.Context, userID int) (string, error)
	AddFarm(с context.Context, farm Farm) error
	DelFarm(с context.Context, userID int, farmID int) error

	AddCow(c context.Context, cow Cow) error
	GetCows(c context.Context, farmID int) (string, error)
	DeleteCows(c context.Context, CowIDs []int) error
	GetCowInfo(c context.Context, farmID int) (string, error)

	GetCowBreeds(c context.Context) (string, error)
	GetBolusesTypes(c context.Context) (string, error)
	AddMonitoringData(c context.Context, data MonitoringData) error
	UpdateHealth(c context.Context, data Health) error
	HasBolus(c context.Context, BolusNum int) int
	GetMonitoringData(c context.Context, cowID int, interval int) ([]MonitoringData, error)
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
