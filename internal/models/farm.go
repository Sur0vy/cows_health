package models

import "context"

type Farm struct {
	ID      int    `json:"farm_id" db:"farm_id,omitempty"`
	Name    string `json:"name" db:"name"`
	Address string `json:"address" db:"address" `
	UserID  int    `json:"-" db:"user_id"`
}

type FarmStorage interface {
	Get(с context.Context, userID int) ([]Farm, error)
	Add(с context.Context, farm Farm) error
	Delete(с context.Context, farmID int) error
}
