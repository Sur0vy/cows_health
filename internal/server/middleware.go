package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Sur0vy/cows_health.git/config"
	"github.com/Sur0vy/cows_health.git/internal/storages"
)

func AuthMiddleware(s storages.UserStorage) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(config.Cookie)
			if err == nil && cookie.Value != "" {
				u := s.Get(c.Request().Context(), cookie.Value)
				if u != nil {
					c.Set("UserID", u.ID)
					return next(c)
				}
			}
			type Unauthorized struct {
				Message string
			}
			msg := &Unauthorized{
				Message: "Unauthorized",
			}
			return c.JSON(http.StatusUnauthorized, msg)
		}
	}
}
