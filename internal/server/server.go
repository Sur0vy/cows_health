package server

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/Sur0vy/cows_health.git/internal/handlers"
	"github.com/Sur0vy/cows_health.git/internal/storage"
)

func SetupServer(s *storage.Storage) *gin.Engine {

	handler := handlers.NewBaseHandler(s)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))

	api := router.Group("/api")

	user := api.Group("/user")
	user.POST("/register", handler.Register)
	user.POST("/login", handler.Login)
	user.POST("/logout", handler.Logout)

	farms := api.Group("/farms")
	farms.Use(CookieMidlewared(s))
	farms.GET("", handler.GetFarms)
	farms.POST("", handler.AddFarm)
	farms.DELETE("/:id", handler.DelFarm)

	//farms.GET("/:id/info", handler.GetFarmInfo)
	//farms.GET("/:id/cows", handler.GetCows)
	//
	//boluses := api.Group("/boluses")
	//boluses.Use(CookieMidlewared(s))
	//boluses.GET("/types", handler.GetBolusesTypes)
	//boluses.POST("/data", handler.AddBolusData)
	//
	cows := api.Group("/cows")
	cows.Use(CookieMidlewared(s))
	//cows.GET("/types", handler.GetCowBreeds)
	//cows.POST("/new", handler.AddCow)
	//cows.DELETE("", handler.DelCows)
	//cows.GET(":id/info", handler.GetCowInfo)
	//cows.PUT(":id/update", CowUpdate)

	//router.GET("/ping", handler.Ping)

	router.NoRoute(handler.ResponseBadRequest)
	return router
}
