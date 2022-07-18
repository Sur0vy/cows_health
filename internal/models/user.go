package models

type User struct {
	ID       int    `json:"-" db:"user_id"`
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
}
