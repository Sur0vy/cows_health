package monitoring_data

import (
	"github.com/Sur0vy/cows_health.git/internal/logger"
	"github.com/Sur0vy/cows_health.git/internal/models"
	"github.com/Sur0vy/cows_health.git/internal/usecase/dataprocessor"
	"github.com/labstack/echo/v4"
)

func Init(group *echo.Group, mt models.MonitoringDataStorage, dp dataprocessor.Processor, log *logger.Logger) {
	mHandler := NewMonitoringDataHandler(mt, dp, log)

	group.POST("", mHandler.Add)
}
