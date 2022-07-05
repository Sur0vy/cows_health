package cow

import (
	"github.com/Sur0vy/cows_health.git/internal/logger"
	"github.com/Sur0vy/cows_health.git/internal/models"
	"github.com/labstack/echo/v4"
)

func Init(group *echo.Group, ct models.CowStorage, log *logger.Logger) {
	cHandler := NewCowHandler(ct, log)

	group.POST("", cHandler.Add)
	group.DELETE("", cHandler.Delete)
	group.GET("/:id", cHandler.Get)
	group.GET("/breeds", cHandler.GetBreeds)
	group.GET("/info/:id", cHandler.GetInfo)
}
