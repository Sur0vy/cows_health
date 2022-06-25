package storage

import "context"

type User struct {
	ID       int    `json:"-"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserStorage interface {
	AddUser(с context.Context, user User) error
	GetUserHash(с context.Context, user User) (string, error)
	GetUser(с context.Context, userHash string) *User
}
