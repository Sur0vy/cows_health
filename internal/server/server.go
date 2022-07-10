package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Sur0vy/cows_health.git/internal/storages"
	"github.com/Sur0vy/cows_health.git/internal/usecase/cow"
	"github.com/Sur0vy/cows_health.git/internal/usecase/dataprocessor"
	"github.com/Sur0vy/cows_health.git/internal/usecase/farm"
	"github.com/Sur0vy/cows_health.git/internal/usecase/monitoringdata"
	"github.com/Sur0vy/cows_health.git/internal/usecase/user"
	"github.com/Sur0vy/cows_health.git/logger"
)

func SetupServer(us storages.UserStorage, fs storages.FarmStorage,
	ms storages.MonitoringDataStorage, cs storages.CowStorage, log *logger.Logger) *echo.Echo {

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

	cowGrp := api.Group("/cow", AuthMiddleware(us))
	cow.Init(cowGrp, cs, log)

	dp := dataprocessor.NewProcessor(ms, cs, log)
	mdGrp := api.Group("/data", AuthMiddleware(us))
	monitoringdata.Init(mdGrp, ms, dp, log)

	return router
}
