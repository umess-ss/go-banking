package transaction

import "time"

type AmountRequest struct {
	Amount float64 `json:"amount"`
}

type TransferRequest struct {
	FromAccountID int64   `json:"from_account_id"`
	ToAccountID   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

type Transaction struct {
	ID              int64     `json:"id"`
	Type            string    `json:"type"`
	FromAccountID   *int64    `json:"from_account_id,omitempty"`
	ToAccountID     *int64    `json:"to_account_id,omitempty"`
	Amount          float64   `json:"amount"`
	Status          string    `json:"status"`
	ReferenceNumber string    `json:"reference_number"`
	CreatedAt       time.Time `json:"created_at"`
}
