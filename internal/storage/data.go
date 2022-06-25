package storage

import (
	"context"
	"time"
)

type Farm struct {
	ID      int    `json:"farm_id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	UserID  int    `json:"-"`
}

type Cow struct {
	ID         int       `json:"id,omitempty"`
	Name       string    `json:"name"`
	BreedID    int       `json:"breed_id,omitempty"`
	Breed      string    `json:"breed,omitempty"`
	FarmID     int       `json:"farm_id,omitempty"`
	BolusNum   int       `json:"bolus_sn"`
	DateOfBorn time.Time `json:"date_of_born"`
	AddedAt    time.Time `json:"added_at"`
	BolusType  string    `json:"bolus_type"`
}

type MonitoringData struct {
	ID          int       `json:"id,omitempty"`
	BolusNum    int       `json:"bolus_sn,omitempty"`
	CowID       int       `json:"cow_id,omitempty"`
	AddedAt     time.Time `json:"added_at"`
	PH          float64   `json:"ph"`
	Temperature float64   `json:"temperature"`
	Movement    float64   `json:"movement"`
	Charge      float64   `json:"charge"`
}

type Health struct {
	CowID     int       `json:"id,omitempty"`
	Ill       string    `json:"ill"`
	Estrus    bool      `json:"estrus"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CowBreed struct {
	ID    int    `json:"breed_id"`
	Breed string `json:"breed"`
}

type CowInfo struct {
	Health  Health           `json:"health"`
	Summary Cow              `json:"summary"`
	History []MonitoringData `json:"history"`
}

type FarmStorage interface {
	GetFarms(с context.Context, userID int) (string, error)
	AddFarm(с context.Context, farm Farm) error
	DelFarm(с context.Context, farmID int) error
	GetCows(c context.Context, farmID int) (string, error)
	AddCow(c context.Context, cow Cow) error
	DeleteCows(c context.Context, CowIDs []int) error
	UpdateHealth(c context.Context, data Health) error
	GetCowInfo(c context.Context, farmID int) (string, error)
	GetCowBreeds(c context.Context) (string, error)
	HasBolus(c context.Context, BolusNum int) int
	GetBolusesTypes(c context.Context) (string, error)
	AddMonitoringData(c context.Context, data MonitoringData) error
	GetMonitoringData(c context.Context, cowID int, interval int) ([]MonitoringData, error)
}
