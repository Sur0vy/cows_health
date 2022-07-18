package storages

import (
	"github.com/Sur0vy/cows_health.git/logger"
	"github.com/jmoiron/sqlx"
)

type StorageDB struct {
	Us *UserStorageDB
	Fs *FarmStorageDB
	Cs *CowStorageDB
	Ms *MonotoringDataStorageDB
}

func NewStorageDB(db *sqlx.DB, log *logger.Logger) *StorageDB {
	return &StorageDB{
		NewUserDB(db, log),
		NewFarmDB(db, log),
		NewCowDB(db, log),
		NewMonitoringDataDB(db, log),
	}
}
