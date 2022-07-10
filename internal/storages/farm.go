package storages

import (
	"context"

	"github.com/Sur0vy/cows_health.git/internal/models"
)

type FarmStorage interface {
	Get(с context.Context, userID int) ([]models.Farm, error)
	Add(с context.Context, farm models.Farm) error
	Delete(с context.Context, farmID int) error
}
