package models

import "time"

type Account struct {
	ID            int64     `json:"id"`
	UserID        *int64    `json:"user_id,omitempty"`
	Name          string    `json:"name"`
	AccountNumber string    `json:"account_number"`
	AccountType   string    `json:"account_type"`
	Balance       float64   `json:"balance"`
	Currency      string    `json:"currency"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
