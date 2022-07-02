package user

import (
	"github.com/Sur0vy/cows_health.git/internal/logger"
	"github.com/Sur0vy/cows_health.git/internal/models"
	"github.com/labstack/echo/v4"
)

func Init(group *echo.Group, st models.UserStorage, log *logger.Logger) {
	uHandler := NewUserHandler(st, log)

	user := group.Group("/user")
	user.POST("/register", uHandler.Register)
	user.POST("/login", uHandler.Login)
	user.POST("/logout", uHandler.Logout)
}
