package server

import (
	"github.com/Sur0vy/cows_health.git/internal/logger"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/Sur0vy/cows_health.git/internal/handlers"
	"github.com/Sur0vy/cows_health.git/internal/storage"
)

func SetupServer(us storage.UserStorage, ds storage.FarmStorage, log *logger.Logger) *gin.Engine {

	handler := handlers.NewBaseHandler(us, ds, log)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))

	api := router.Group("/api")

	user := api.Group("/user")
	user.POST("/register", handler.Register)
	user.POST("/login", handler.Login)
	user.POST("/logout", handler.Logout)

	farms := api.Group("/farms")
	farms.Use(CookieMidlewared(us))
	farms.GET("", handler.GetFarms)
	farms.POST("", handler.AddFarm)
	farms.DELETE("/:id", handler.DelFarm)
	farms.GET("/:id/cows", handler.GetCows)

	boluses := api.Group("/boluses")
	boluses.Use(CookieMidlewared(us))
	boluses.GET("/types", handler.GetBolusesTypes)
	boluses.POST("/data", handler.AddMonitoringData)

	cows := api.Group("/cows")
	cows.Use(CookieMidlewared(us))
	cows.GET("/breeds", handler.GetCowBreeds)
	cows.POST("", handler.AddCow)
	cows.DELETE("", handler.DelCows)
	cows.GET(":id/info", handler.GetCowInfo)

	router.NoRoute(handler.ResponseBadRequest)
	return router
}
