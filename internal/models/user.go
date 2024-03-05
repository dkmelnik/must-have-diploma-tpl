package models

import "time"

type ModelID string

type User struct {
	ID        ModelID   `db:"id"`
	Login     string    `db:"login"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}
