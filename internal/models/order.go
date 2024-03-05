package models

import (
	"database/sql"
	"time"
)

type OrderStatus string

var (
	OrderNew        OrderStatus = "NEW"
	OrderProcessing OrderStatus = "PROCESSING"
	OrderRegistered OrderStatus = "REGISTERED"
	OrderInvalid    OrderStatus = "INVALID"
	OrderProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID        ModelID         `db:"id"`
	UserID    ModelID         `db:"user_id"`
	Number    string          `db:"number"`
	Accrual   sql.NullFloat64 `db:"accrual"`
	Status    OrderStatus     `db:"status"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}

func (m *Order) SetAccrual(accrual *float64) {
	if accrual != nil {
		m.Accrual = sql.NullFloat64{
			Float64: *accrual,
			Valid:   true,
		}
	}
}
