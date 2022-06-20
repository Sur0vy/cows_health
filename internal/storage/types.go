package storage

import (
	"time"
)

type User struct {
	ID       int    `json:"-"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Farm struct {
	ID      int    `json:"farm_id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	UserID  int    `json:"-"`
}

type CowBreed struct {
	ID    int    `json:"breed_id"`
	Breed string `json:"breed"`
}

type Cow struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	BreedID    int       `json:"breed_id"`
	FarmID     int       `json:"farm_id,omitempty"`
	BolusNum   int       `json:"bolus_sn"`
	DateOfBorn time.Time `json:"date_of_born"`
	AddedAt    time.Time `json:"added_at"`
	BolusType  string    `json:"bolus_type"`
}

type MonitoringData struct {
	ID          int       `json:"id"`
	BolusNum    int       `json:"bolus_sn"`
	CowID       int       `json:"cow_id"`
	AddedAt     time.Time `json:"added_at"`
	PH          float32   `json:"ph"`
	Temperature float32   `json:"temperature"`
	Movement    float32   `json:"movement"`
	Charge      float32   `json:"charge"`
}

type Health struct {
	CowID     int       `json:"id"`
	Ill       string    `json:"ill"`
	Estrus    bool      `json:"estrus"`
	UpdatedAt time.Time `json:"updated_at"`
}

//
//type CowInfo struct {
//	Health  Health           `json:"health"`
//	Summary Cow              `json:"summary"`
//	History []MonitoringData `json:"history"`
//}
