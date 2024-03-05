package models

import (
	"time"
)

type Withdrawal struct {
	ID          ModelID   `db:"id"`
	UserID      ModelID   `db:"user_id"`
	OrderNumber string    `db:"order_number"`
	Amount      float64   `db:"amount"`
	CreatedAt   time.Time `db:"created_at"`
}
