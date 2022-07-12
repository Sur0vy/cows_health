package cow

import (
	"github.com/labstack/echo/v4"

	"github.com/Sur0vy/cows_health.git/internal/storages"
	"github.com/Sur0vy/cows_health.git/logger"
)

func Init(group *echo.Group, ct storages.CowStorage, log *logger.Logger) {
	cHandler := NewCowHandler(ct, log)

	group.POST("", cHandler.Add)
	group.DELETE("", cHandler.Delete)
	group.GET("/:id", cHandler.Get)
	group.GET("/breeds", cHandler.GetBreeds)
	group.POST("/breed", cHandler.AddBreed)
	group.GET("/info/:id", cHandler.GetInfo)
}
