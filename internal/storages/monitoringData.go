package storages

import (
	"context"

	"github.com/Sur0vy/cows_health.git/internal/models"
)

type MonitoringDataStorage interface {
	Add(c context.Context, data models.MonitoringData) error
	Get(c context.Context, cowID int, interval int) ([]models.MonitoringData, error)
}
