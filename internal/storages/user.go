package storages

import (
	"context"

	"github.com/Sur0vy/cows_health.git/internal/models"
)

type UserStorage interface {
	Add(с context.Context, user models.User) error
	GetHash(с context.Context, user models.User) (string, error)
	Get(с context.Context, userHash string) *models.User
}
