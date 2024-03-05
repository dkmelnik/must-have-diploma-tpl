package dto

import "time"

type (
	WithdrawalPayload struct {
		Order string  `json:"order"`
		Sum   float64 `json:"sum"`
	}
	WithdrawalResponse struct {
		Order       string    `json:"order"`
		Sum         float64   `json:"sum"`
		ProcessedAT time.Time `json:"processed_at"`
	}
)
