package storages

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"

	"github.com/Sur0vy/cows_health.git/internal/errors"
	"github.com/Sur0vy/cows_health.git/internal/helpers"
	"github.com/Sur0vy/cows_health.git/internal/logger"
	"github.com/Sur0vy/cows_health.git/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserStorageDB struct {
	log *logger.Logger
	db  *sqlx.DB
}

func NewUserDB(db *sqlx.DB, log *logger.Logger) *UserStorageDB {
	return &UserStorageDB{
		log: log,
		db:  db,
	}
}

func (s *UserStorageDB) Add(c context.Context, user models.User) error {
	userHash := helpers.GetMD5Hash(user.Login)
	u := s.Get(c, userHash)
	if u != nil {
		s.log.Warn().Msgf("User %v already exists", user.Login)
		return errors.NewExistError()
	}

	passwordHash, err := helpers.GetCryptoPassword(user.Password)
	if err != nil {
		s.log.Warn().Msg("Error than encrypting password")
		return err
	}

	ctxIn, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	dialect := goqu.Dialect("postgres")
	//добавление
	ds := dialect.Insert("users").
		Cols("login", "password").
		Vals(
			goqu.Vals{userHash, passwordHash},
		)
	SQLStr, _, _ := ds.ToSQL()

	_, err = s.db.ExecContext(ctxIn, SQLStr)
	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return err
	}
	return nil
}

func (s *UserStorageDB) Get(c context.Context, userHash string) *models.User {
	ctxIn, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	usr := models.User{}
	dialect := goqu.Dialect("postgres")
	SQLStr, _, _ := dialect.From("users").
		Where(
			goqu.C("login").Eq(userHash),
		).ToSQL()

	err := s.db.GetContext(ctxIn, &usr, SQLStr)
	if err != nil {
		s.log.Warn().Err(err).Msg("db request error")
		return nil
	}
	return &usr
}

func (s *UserStorageDB) GetHash(c context.Context, user models.User) (string, error) {
	userHash := helpers.GetMD5Hash(user.Login)
	u := s.Get(c, userHash)
	if u == nil {
		s.log.Warn().Msgf("User %v not exists", user.Login)
		return "", errors.NewEmptyError()
	}
	if helpers.CheckPassword(u.Password, user.Password) {
		return userHash, nil
	}
	return "", errors.NewEmptyError()
}
