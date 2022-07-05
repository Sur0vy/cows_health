package models

import (
	"context"
	"time"
)

type Breed struct {
	ID   int    `json:"breed_id" db:"breed_id"`
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

type CowStorage interface {
	Add(c context.Context, cow Cow) error
	Get(c context.Context, farmID int) ([]Cow, error)
	Delete(c context.Context, CowIDs []int) error
	GetInfo(c context.Context, farmID int) (CowInfo, error)
	GetBreeds(c context.Context) ([]Breed, error)
	UpdateHealth(c context.Context, data Health) error
	HasBolus(c context.Context, BolusNum int) int
}
