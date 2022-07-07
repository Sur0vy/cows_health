package user

import (
	"github.com/Sur0vy/cows_health.git/internal/models"
	"github.com/Sur0vy/cows_health.git/logger"
	"github.com/labstack/echo/v4"
)

func Init(group *echo.Group, st models.UserStorage, log *logger.Logger) {
	uHandler := NewUserHandler(st, log)

	group.POST("/register", uHandler.Register)
	group.POST("/login", uHandler.Login)
	group.POST("/logout", uHandler.Logout)
}
