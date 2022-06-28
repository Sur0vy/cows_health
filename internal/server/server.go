package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Sur0vy/cows_health.git/internal/handlers"
	"github.com/Sur0vy/cows_health.git/internal/logger"
	"github.com/Sur0vy/cows_health.git/internal/storage"
)

func SetupServer(us storage.UserStorage, fs storage.FarmStorage, log *logger.Logger) *echo.Echo {

	uHandler := handlers.NewUserHandler(us, log)
	fHandler := handlers.NewFarmHandler(fs, log)

	router := echo.New()

	router.Use(middleware.Gzip())

	router.Any("/*", uHandler.ResponseBadRequest)

	api := router.Group("/api")

	user := api.Group("/user")
	user.POST("/register", uHandler.Register)
	user.POST("/login", uHandler.Login)
	user.POST("/logout", uHandler.Logout)

	farms := api.Group("/farms", AuthMiddleware(us))
	farms.GET("", fHandler.GetFarms)
	farms.POST("", fHandler.AddFarm)
	farms.DELETE("/:id", fHandler.DelFarm)
	farms.GET("/:id/cows", fHandler.GetCows)

	boluses := api.Group("/boluses", AuthMiddleware(us))
	boluses.GET("/types", fHandler.GetBolusesTypes)
	boluses.POST("/data", fHandler.AddMonitoringData)

	cows := api.Group("/cows", AuthMiddleware(us))
	cows.GET("/breeds", fHandler.GetCowBreeds)
	cows.POST("", fHandler.AddCow)
	cows.DELETE("", fHandler.DelCows)
	cows.GET(":id/info", fHandler.GetCowInfo)
	return router
}
