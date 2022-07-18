package farm

import (
	"github.com/labstack/echo/v4"

	"github.com/Sur0vy/cows_health.git/internal/storages"
	"github.com/Sur0vy/cows_health.git/logger"
)

func Init(group *echo.Group, st storages.FarmStorage, log *logger.Logger) {
	fHandler := NewFarmHandler(st, log)

	group.GET("", fHandler.Get)
	group.POST("", fHandler.Add)
	group.DELETE("/:id", fHandler.Delete)
}
