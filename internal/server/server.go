package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Sur0vy/cows_health.git/internal/handlers"
	"github.com/Sur0vy/cows_health.git/internal/logger"
	"github.com/Sur0vy/cows_health.git/internal/storage"
)

func SetupServer(us storage.UserStorage, ds storage.FarmStorage, log *logger.Logger) *echo.Echo {

	handler := handlers.NewBaseHandler(us, ds, log)
	router := echo.New()

	router.Use(middleware.Gzip())
	
	router.Any("/*", handler.ResponseBadRequest)

	api := router.Group("/api")

	user := api.Group("/user")
	user.POST("/register", handler.Register)
	user.POST("/login", handler.Login)
	user.POST("/logout", handler.Logout)

	farms := api.Group("/farms", AuthMiddleware(us))
	farms.GET("", handler.GetFarms)
	farms.POST("", handler.AddFarm)
	farms.DELETE("/:id", handler.DelFarm)
	farms.GET("/:id/cows", handler.GetCows)

	boluses := api.Group("/boluses", AuthMiddleware(us))
	boluses.GET("/types", handler.GetBolusesTypes)
	boluses.POST("/data", handler.AddMonitoringData)

	cows := api.Group("/cows", AuthMiddleware(us))
	cows.GET("/breeds", handler.GetCowBreeds)
	cows.POST("", handler.AddCow)
	cows.DELETE("", handler.DelCows)
	cows.GET(":id/info", handler.GetCowInfo)
	return router
}
