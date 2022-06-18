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
	FarmID     int       `json:"farm_id"`
	BolusNum   int       `json:"bolus_sn"`
	DateOfBorn time.Time `json:"date_of_born"`
	AddedAt    time.Time `json:"added_at"`
	BolusType  string    `json:"bolus_type"`
}

//type Health struct {
//	Drink       bool    `json:"drink"`
//	Stress      bool    `json:"stress"`
//	Temperature float32 `json:"temperature"`
//	Activity    float32 `json:"activity"`
//	CowID       int     `json:"-"`
//}
//
//type MonitoringData struct {
//	BolusID      int     `json:"-"`
//	SerialNumber string  `json:"num"`
//	Time         string  `json:"dateTime"`
//	PH           float32 `json:"ph"`
//	Temperature  float32 `json:"temperature"`
//	Movement     float32 `json:"movement"`
//	Humidity     float32 `json:"humidity"`
//	Charge       float32 `json:"charge"`
//}
//
//type CowInfo struct {
//	Health  Health           `json:"health"`
//	Summary Cow              `json:"summary"`
//	History []MonitoringData `json:"history"`
//}
//
//type BolusType struct {
//	ID   int
//	Name string
//}
