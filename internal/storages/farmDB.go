package storages

import (
	"context"
	"database/sql"
	"github.com/Sur0vy/cows_health.git/errors"
	"github.com/Sur0vy/cows_health.git/logger"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"

	"github.com/Sur0vy/cows_health.git/internal/models"
)

type FarmStorageDB struct {
	log *logger.Logger
	db  *sqlx.DB
}

func NewFarmDB(db *sqlx.DB, log *logger.Logger) *FarmStorageDB {
	return &FarmStorageDB{
		log: log,
		db:  db,
	}
}

func (s *FarmStorageDB) Get(c context.Context, userID int) ([]models.Farm, error) {
	ctxIn, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var farms []models.Farm

	dialect := goqu.Dialect("postgres")
	SQLStr, _, _ := dialect.From("farms").
		Select("farm_id", "name", "address").
		Where(
			goqu.C("user_id").Eq(userID),
			goqu.C("deleted").Neq(true),
		).ToSQL()

	err := s.db.SelectContext(ctxIn, &farms, SQLStr)
	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return farms, err
	}

	if len(farms) == 0 {
		return farms, errors.ErrEmpty
	}
	return farms, nil
}

func (s *FarmStorageDB) Add(c context.Context, farm models.Farm) error {
	ctxIn, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	dialect := goqu.Dialect("postgres")
	SQLStr, _, _ := dialect.
		Select("farm_id").
		From("farms").
		Where(
			goqu.C("address").Eq(farm.Address),
		).ToSQL()
	var farmID int
	err := s.db.GetContext(ctxIn, &farmID, SQLStr)

	if err == nil {
		s.log.Info().Msg("farm already exist")
		return errors.ErrExist
	} else if err != sql.ErrNoRows {
		s.log.Warn().Err(err).Msg("db request error")
		return err
	}

	//добавление
	ds := dialect.
		Insert("farms").
		Cols("name", "address", "user_id").
		Vals(
			goqu.Vals{farm.Name, farm.Address, farm.UserID},
		)
	SQLStr, _, _ = ds.ToSQL()

	_, err = s.db.ExecContext(ctxIn, SQLStr)
	if err != nil {
		s.log.Warn().Err(err).Msg("inserting farm error")
		return err
	}
	return nil
}

func (s *FarmStorageDB) Delete(c context.Context, farmID int) error {
	ctxIn, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	dialect := goqu.Dialect("postgres")
	SQLStrF, _, _ := dialect.Update("farms").
		Set(
			goqu.Record{"deleted": "true"},
		).
		Where(
			goqu.C("farm_id").Eq(farmID),
			goqu.C("deleted").Neq(true),
		).ToSQL()
	SQLStrC, _, _ := dialect.Update("cows").
		Set(
			goqu.Record{"deleted": "true"},
		).
		Where(
			goqu.C("farm_id").Eq(farmID),
			goqu.C("deleted").Neq(true),
		).ToSQL()

	tx := s.db.MustBegin()
	//фермы
	tx.MustExecContext(ctxIn, SQLStrF)
	//коровы
	tx.MustExecContext(ctxIn, SQLStrC)
	err := tx.Commit()
	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return err
	}
	return nil
}
