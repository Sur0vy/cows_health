package models

import (
	"context"
	"time"
)

type MonitoringData struct {
	ID          int       `json:"id,omitempty" db:"md_id"`
	BolusNum    int       `json:"bolus_sn,omitempty"`
	CowID       int       `json:"cow_id,omitempty" db:"cov_id"`
	AddedAt     time.Time `json:"added_at" db:"added_at"`
	PH          float64   `json:"ph" db:"ph"`
	Temperature float64   `json:"temperature" db:"temperature"`
	Movement    float64   `json:"movement" db:"movement"`
	Charge      float64   `json:"charge" db:"charge"`
}

type MonitoringDataStorage interface {
	Add(c context.Context, data MonitoringData) error
	Get(c context.Context, cowID int, interval int) ([]MonitoringData, error)
}
