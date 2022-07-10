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

type CowStorageDB struct {
	log *logger.Logger
	db  *sqlx.DB
}

func NewCowDB(db *sqlx.DB, log *logger.Logger) *CowStorageDB {
	return &CowStorageDB{
		log: log,
		db:  db,
	}
}

func (s *CowStorageDB) Add(c context.Context, cow models.Cow) error {
	cowID := s.HasBolus(c, cow.BolusNum)

	if cowID != -1 {
		s.log.Info().Msg("duplicate bolus")
		return errors.ErrExist
	}

	ctxIn, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	dialect := goqu.Dialect("postgres")
	SQLStr, _, _ := dialect.
		Insert("cows").
		Cols("name", "breed_id", "farm_id", "bolus_sn", "date_of_born", "added_at").
		Vals(
			goqu.Vals{cow.Name, cow.BreedID, cow.FarmID, cow.BolusNum, cow.DateOfBorn, cow.AddedAt}).
		Returning("cow_id").
		ToSQL()

	err := s.db.GetContext(ctxIn, &cow.ID, SQLStr)
	if err != nil {
		s.log.Info().Msg("inserting cow error")
		return err
	}

	SQLStr, _, _ = dialect.
		Insert("health").
		Cols("cow_id").
		Vals(
			goqu.Vals{cow.ID},
		).ToSQL()

	_, err = s.db.ExecContext(ctxIn, SQLStr)
	if err != nil {
		s.log.Warn().Err(err).Msg("inserting cow error")
		return err
	}
	return nil
}

func (s *CowStorageDB) Delete(c context.Context, CowIDs []int) error {
	ctxIn, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	dialect := goqu.Dialect("postgres")
	SQLStr, _, _ := dialect.Update("cows").
		Set(
			goqu.Record{"deleted": "true"},
		).
		Where(
			goqu.C("cow_id").In(CowIDs),
			goqu.C("deleted").Neq(true),
		).ToSQL()
	_, err := s.db.ExecContext(ctxIn, SQLStr)
	if err != nil {
		s.log.Warn().Err(err).Msg("deleting cows error")
		return err
	}
	return nil
}

func (s *CowStorageDB) GetBreeds(c context.Context) ([]models.Breed, error) {
	ctxIn, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	var breeds []models.Breed

	dialect := goqu.Dialect("postgres")
	SQLStr, _, _ := dialect.From("breeds").ToSQL()

	err := s.db.SelectContext(ctxIn, &breeds, SQLStr)
	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return breeds, err
	}
	return breeds, nil
}

func (s *CowStorageDB) Get(c context.Context, farmID int) ([]models.Cow, error) {
	ctxIn, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	var cows []models.Cow
	dialect := goqu.Dialect("postgres")
	SQLStr, _, _ := dialect.From("cows").
		Select("cow_id", "name", "breed_id", "farm_id", "bolus_sn", "date_of_born", "added_at").
		Where(
			goqu.C("farm_id").Eq(farmID),
			goqu.C("deleted").Neq(true),
		).ToSQL()

	err := s.db.SelectContext(ctxIn, &cows, SQLStr)
	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return cows, err
	}

	if len(cows) == 0 {
		return cows, errors.ErrEmpty
	}
	return cows, nil
}

func (s *CowStorageDB) GetInfo(c context.Context, cowID int) (models.CowInfo, error) {
	ctxIn, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	hour := 3600
	intervalInS := hour * 24

	var cowInfo models.CowInfo

	dialect := goqu.Dialect("postgres")
	SQLStr, _, _ := dialect.From("cows").
		Select("cows.name", "breeds.name", "cows.bolus_sn", "cows.date_of_born",
			"health.estrus", "health.ill", "health.updated_at",
			"monitoring_data.added_at", "monitoring_data.ph", "monitoring_data.temperature",
			"monitoring_data.movement", "monitoring_data.charge").
		Join(
			goqu.T("breeds"),
			goqu.On(goqu.Ex{"cows.breed_id": goqu.I("breeds.breed_id")}),
		).
		Join(
			goqu.T("health"),
			goqu.On(goqu.Ex{"cows.cow_id": goqu.I("health.cow_id")}),
		).
		Join(
			goqu.T("monitoring_data"),
			goqu.On(goqu.Ex{"cows.cow_id": goqu.I("monitoring_data.cow_id")}),
		).
		Where(
			//goqu.Ex{"cows.cow_id": goqu.(cowID)},
			goqu.I("cows.cow_id").Eq(cowID),
			goqu.V("EXTRACT(EPOCH FROM now()) - EXTRACT(EPOCH FROM monitoring_data.added_at)").Lt(intervalInS),
		).
		Order(
			goqu.I("cows.added_at").Asc()).
		ToSQL()
	var re = regexp.MustCompile(`'`)
	SQLStr = re.ReplaceAllString(SQLStr, "")

	rows, err := s.db.QueryxContext(ctxIn, SQLStr)
	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return cowInfo, err
	}
	defer rows.Close()

	for rows.Next() {
		var md models.MonitoringData
		err = rows.Scan(&cowInfo.Summary.Name, &cowInfo.Summary.Breed, &cowInfo.Summary.BolusNum,
			&cowInfo.Summary.DateOfBorn,
			&cowInfo.Health.Estrus, &cowInfo.Health.Ill, &cowInfo.Health.UpdatedAt,
			&md.AddedAt, &md.PH, &md.Temperature, &md.Movement, &md.Charge)
		if err != nil {
			s.log.Warn().Err(err).Msg("get monitoring data instance error")
			return cowInfo, err
		}
		cowInfo.History = append(cowInfo.History, md)
	}

	if err := rows.Err(); err != nil {
		s.log.Warn().Err(err).Msg("get cow info rows error")
		return cowInfo, err
	}

	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return cowInfo, err
	}

	return cowInfo, nil
}

func (s *CowStorageDB) HasBolus(c context.Context, BolusNum int) int {
	ctxIn, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	var cowID int

	dialect := goqu.Dialect("postgres")
	SQLStr, _, _ := dialect.From("cows").
		Select("cow_id").
		Limit(1).
		Where(
			goqu.C("bolus_sn").Eq(BolusNum),
		).ToSQL()

	err := s.db.GetContext(ctxIn, &cowID, SQLStr)
	if err == nil {
		return cowID
	}
	return -1
}

func (s *CowStorageDB) UpdateHealth(c context.Context, data models.Health) error {
	ctxIn, cancel := context.WithTimeout(c, time.Second)
	defer cancel()

	dialect := goqu.Dialect("postgres")
	SQLStr, _, _ := dialect.Update("health").
		Set(
			goqu.Record{
				"updated_at": data.UpdatedAt,
				"ill":        data.Ill,
				"estrus":     data.Estrus,
			},
		).
		Where(
			goqu.C("cow_id").Eq(data.CowID),
		).
		ToSQL()

	_, err := s.db.ExecContext(ctxIn, SQLStr)
	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return err
	}
	return nil
}

//func (s *CowStorageDB) Ping() error {
//	return s.db.Ping()
//}
