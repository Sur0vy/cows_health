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

func SetupServer(s storages.StorageDB, log *logger.Logger) *echo.Echo {

	router := echo.New()
	router.Use(middleware.Gzip())
	//any rout
	router.Any("/*", func(c echo.Context) error {
		log.Info().Msgf("bad request. Error code: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	})
	api := router.Group("/api")
	userGrp := api.Group("/user")
	user.Init(userGrp, s.Us, log)

	farmGrp := api.Group("/farm", AuthMiddleware(s.Us))
	farm.Init(farmGrp, s.Fs, log)

	cowGrp := api.Group("/cow", AuthMiddleware(s.Us))
	cow.Init(cowGrp, s.Cs, log)

	dp := dataprocessor.NewProcessor(s.Ms, s.Cs, log)
	mdGrp := api.Group("/data", AuthMiddleware(s.Us))
	monitoringdata.Init(mdGrp, s.Ms, dp, log)

	return router
}
