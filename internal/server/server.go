package server

import (
	"github.com/Sur0vy/cows_health.git/internal/handlers/farm"
	"github.com/Sur0vy/cows_health.git/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	//"github.com/Sur0vy/cows_health.git/internal/entity/farm"
	//"github.com/Sur0vy/cows_health.git/internal/entity/user"
	"github.com/Sur0vy/cows_health.git/internal/handlers/user"
	"github.com/Sur0vy/cows_health.git/internal/logger"
)

func SetupServer(us models.UserStorage, fs models.FarmStorage, log *logger.Logger) *echo.Echo {
	//fHandler := farm.NewFarmHandler(fs, log)

	router := echo.New()
	router.Use(middleware.Gzip())
	//any rout
	router.Any("/*", func(c echo.Context) error {
		log.Info().Msgf("bad request. Error code: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	})
	api := router.Group("/api")
	userGrp := api.Group("/user")
	user.Init(userGrp, us, log)

	farmGrp := api.Group("/farm", AuthMiddleware(us))
	farm.Init(farmGrp, fs, log)

	//farms.GET("/:id/cows", fHandler.GetCows)
	//
	//boluses := api.Group("/boluses", AuthMiddleware(us))
	//boluses.GET("/types", fHandler.GetBolusesTypes)
	//boluses.POST("/data", fHandler.AddMonitoringData)
	//
	//cows := api.Group("/cows", AuthMiddleware(us))
	//cows.GET("/breeds", fHandler.GetCowBreeds)
	//cows.POST("", fHandler.AddCow)
	//cows.DELETE("", fHandler.DelCows)
	//cows.GET(":id/info", fHandler.GetCowInfo)
	return router
}
