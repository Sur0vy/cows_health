package models

import (
	"time"
)

type MonitoringData struct {
	ID          int       `json:"id,omitempty" db:"-"`
	BolusNum    int       `json:"bolus_sn,omitempty" db:"-"`
	CowID       int       `json:"cow_id,omitempty" db:"cow_id"`
	AddedAt     time.Time `json:"added_at" db:"added_at"`
	PH          float64   `json:"ph" db:"ph"`
	Temperature float64   `json:"temperature" db:"temperature"`
	Movement    float64   `json:"movement" db:"movement"`
	Charge      float64   `json:"charge" db:"charge"`
}

type MonitoringDataFull struct {
	Data        []MonitoringData
	AvgPH       float64
	AvgTemp     float64
	AvgMovement float64
}
