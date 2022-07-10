package models

type Farm struct {
	ID      int    `json:"farm_id" db:"farm_id,omitempty"`
	Name    string `json:"name" db:"name"`
	Address string `json:"address" db:"address" `
	UserID  int    `json:"-" db:"user_id"`
}
