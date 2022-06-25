package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Sur0vy/cows_health.git/internal/config"
	"github.com/Sur0vy/cows_health.git/internal/storage"
)

func CookieMidlewared(s storage.UserStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(config.Cookie)
		if err == nil && cookie != "" {
			u := s.GetUser(context.Background(), cookie)
			if u != nil {
				//				logger.Wr.Info().Msg("Cookie accepted")
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
		c.Abort()
	}
}
