package storages

import (
	"context"

	"github.com/Sur0vy/cows_health.git/internal/models"
)

type CowStorage interface {
	Add(c context.Context, cow models.Cow) error
	Get(c context.Context, farmID int) ([]models.Cow, error)
	Delete(c context.Context, CowIDs []int) error
	GetInfo(c context.Context, farmID int) (models.CowInfo, error)
	GetBreeds(c context.Context) ([]models.Breed, error)
	UpdateHealth(c context.Context, data models.Health) error
	HasBolus(c context.Context, BolusNum int) int
}
