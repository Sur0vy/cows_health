package monitoringdata

import (
	"github.com/labstack/echo/v4"

	"github.com/Sur0vy/cows_health.git/internal/storages"
	"github.com/Sur0vy/cows_health.git/internal/usecase/dataprocessor"
	"github.com/Sur0vy/cows_health.git/logger"
)

func Init(group *echo.Group, mt storages.MonitoringDataStorage, dp dataprocessor.DataProcessor, log *logger.Logger) {
	mHandler := NewMonitoringDataHandler(mt, dp, log)

	group.POST("", mHandler.Add)
}
