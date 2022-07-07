package monitoringdata

import (
	"github.com/Sur0vy/cows_health.git/internal/models"
	"github.com/Sur0vy/cows_health.git/logger"
	"github.com/labstack/echo/v4"
)

func Init(group *echo.Group, mt models.MonitoringDataStorage, dp models.DataProcessor, log *logger.Logger) {
	mHandler := NewMonitoringDataHandler(mt, dp, log)

	group.POST("", mHandler.Add)
}
