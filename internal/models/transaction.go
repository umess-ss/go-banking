package models

import "time"

type AmountRequest struct {
	Amount float64 `json:"amount"`
}

type TransferRequest struct {
	FromAccountID int     `json:"from_account_id"`
	ToAccountID   int     `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

type Transaction struct {
	ID            int       `json:"id"`
	Type          string    `json:"type"`
	FromAccountID *int      `json:"from_account_id,omitempty"`
	ToAccountID   *int      `json:"to_account_id,omitempty"`
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}
