package dataprocessor

import (
	"context"
	"sync"

	"github.com/Sur0vy/cows_health.git/internal/models"
)

type DataProcessor interface {
	Run(c context.Context, data models.MonitoringData, wg *sync.WaitGroup)
	GetHealthData(c context.Context, cowID int) (models.MonitoringDataFull, error)
	CalculateHealth(mds models.MonitoringDataFull) models.Health
	Save(c context.Context, data models.MonitoringData) error
}
