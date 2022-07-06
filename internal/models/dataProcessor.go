package models

import (
	"context"
	"sync"
)

type DataSlice struct {
	Data        []MonitoringData
	AvgPH       float64
	AvgTemp     float64
	AvgMovement float64
}

type DataProcessor interface {
	Run(c context.Context, data MonitoringData, wg *sync.WaitGroup)
	GetHealthData(c context.Context, cowID int) (DataSlice, error)
	CalculateHealth(mds DataSlice) Health
	Save(c context.Context, data MonitoringData) error
}
