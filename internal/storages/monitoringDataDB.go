package storages

import (
	"context"
	"regexp"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"

	"github.com/Sur0vy/cows_health.git/errors"
	"github.com/Sur0vy/cows_health.git/internal/models"
	"github.com/Sur0vy/cows_health.git/logger"
)

type MonotoringDataStorageDB struct {
	log *logger.Logger
	db  *sqlx.DB
}

func NewMonitoringDataDB(db *sqlx.DB, log *logger.Logger) *MonotoringDataStorageDB {
	return &MonotoringDataStorageDB{
		log: log,
		db:  db,
	}
}

func (s *MonotoringDataStorageDB) Add(c context.Context, data models.MonitoringData) error {
	ctxIn, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	dialect := goqu.Dialect("postgres")

	SQLStr, _, _ := dialect.
		Insert("monitoring_data").
		Rows(data).ToSQL()

	_, err := s.db.ExecContext(ctxIn, SQLStr)
	if err != nil {
		s.log.Warn().Err(err).Msg("inserting farm error")
		return err
	}
	return nil
}

func (s *MonotoringDataStorageDB) Get(c context.Context, cowID int, interval int) ([]models.MonitoringData, error) {
	ctxIn, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	var res []models.MonitoringData

	min := 60
	intervalInS := min * interval

	dialect := goqu.Dialect("postgres")
	SQLStr, _, _ := dialect.From("monitoring_data").
		Select("added_at", "ph", "temperature", "movement").
		Where(
			goqu.C("cow_id").Eq(cowID),
			goqu.V("EXTRACT(EPOCH FROM now()) - EXTRACT(EPOCH FROM added_at)").Lt(intervalInS),
		).ToSQL()

	var re = regexp.MustCompile(`'`)
	SQLStr = re.ReplaceAllString(SQLStr, "")

	err := s.db.SelectContext(ctxIn, &res, SQLStr)
	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return res, err
	}

	if len(res) == 0 {
		return res, errors.ErrEmpty
	}
	return res, nil
}
