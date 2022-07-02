package models

import "context"

type User struct {
	ID       int    `json:"-" db:"user_id"`
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
}

type UserStorage interface {
	Add(с context.Context, user User) error
	GetHash(с context.Context, user User) (string, error)
	Get(с context.Context, userHash string) *User
}
