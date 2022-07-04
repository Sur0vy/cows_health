package farm

import (
	"github.com/Sur0vy/cows_health.git/internal/logger"
	"github.com/Sur0vy/cows_health.git/internal/models"
	"github.com/labstack/echo/v4"
)

func Init(group *echo.Group, st models.FarmStorage, log *logger.Logger) {
	fHandler := NewFarmHandler(st, log)

	group.GET("", fHandler.Get)
	group.POST("", fHandler.Add)
	group.DELETE("/:id", fHandler.Delete)
}
