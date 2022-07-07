package farm

import (
	"github.com/Sur0vy/cows_health.git/logger"
	"github.com/labstack/echo/v4"

	"github.com/Sur0vy/cows_health.git/internal/models"
)

func Init(group *echo.Group, st models.FarmStorage, log *logger.Logger) {
	fHandler := NewFarmHandler(st, log)

	group.GET("", fHandler.Get)
	group.POST("", fHandler.Add)
	group.DELETE("/:id", fHandler.Delete)
}
