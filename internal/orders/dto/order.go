package dto

import "time"

type OrderResponse struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    *float64  `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func (o *OrderResponse) SetAccrual(accrual float64) {
	o.Accrual = &accrual
}
