package models

import (
	"time"
)

//TODO напрашивается вынести в другой файл с созданием отдельного стораджа
type Breed struct {
	ID   int    `json:"breed_id,omitempty" db:"breed_id"`
	Name string `json:"breed" db:"name"`
}

type Cow struct {
	ID         int       `json:"id,omitempty"  db:"cow_id"`
	Name       string    `json:"name" db:"name"`
	BreedID    int       `json:"breed_id,omitempty" db:"breed_id"`
	Breed      string    `json:"breed,omitempty" db:"-"`
	FarmID     int       `json:"farm_id,omitempty" db:"farm_id"`
	BolusNum   int       `json:"bolus_sn" db:"bolus_sn"`
	DateOfBorn time.Time `json:"date_of_born" db:"date_of_born"`
	AddedAt    time.Time `json:"added_at" db:"added_at"`
}

type Health struct {
	CowID     int       `json:"id,omitempty"`
	Ill       string    `json:"ill"`
	Estrus    bool      `json:"estrus"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CowInfo struct {
	Health  Health           `json:"health"`
	Summary Cow              `json:"summary"`
	History []MonitoringData `json:"history"`
}
