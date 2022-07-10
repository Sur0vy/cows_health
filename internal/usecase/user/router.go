package user

import (
	"github.com/labstack/echo/v4"

	"github.com/Sur0vy/cows_health.git/internal/storages"
	"github.com/Sur0vy/cows_health.git/logger"
)

func Init(group *echo.Group, st storages.UserStorage, log *logger.Logger) {
	uHandler := NewUserHandler(st, log)

	group.POST("/register", uHandler.Register)
	group.POST("/login", uHandler.Login)
	group.POST("/logout", uHandler.Logout)
}
