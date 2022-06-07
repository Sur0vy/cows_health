package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Sur0vy/cows_health.git/internal/config"
	"github.com/Sur0vy/cows_health.git/internal/logger"
	"github.com/Sur0vy/cows_health.git/internal/storage"
)

func CookieMidlewared(s *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(config.Cookie)
		logger.Wr.Info().Msgf("Extract cookie: %v", cookie)
		if err == nil && cookie != "" {
			_, err := (*s).GetUser(context.Background(), cookie)
			if err == nil {
				logger.Wr.Info().Msg("Cookie accepted")
				return
			}
		}
		logger.Wr.Info().Msg("Bad cookie or cookie not found")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
		c.Abort()
	}
}
