package repository

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"go-banking/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

//AccountRepository handles data storage and retrieval for accounts.
//In a real application, this would interface with a database, but for simplicity,
// we'll use an in-memory slice.

type AccountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) FindAll(ctx context.Context) ([]models.Account, error) {
	query := `
		SELECT id, user_id, name, account_number, balance, currency, created_at, updated_at
		FROM accounts
		ORDER BY id DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []models.Account{}
	for rows.Next() {
		var account models.Account
		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Name,
			&account.AccountNumber,
			&account.Balance,
			&account.Currency,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (r *AccountRepository) FindByID(ctx context.Context, id int64) (*models.Account, error) {
	query := `
		SELECT id, user_id, name, account_number, balance, currency, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`

	var account models.Account
	err := r.db.QueryRow(ctx, query, id).Scan(
		&account.ID,
		&account.UserID,
		&account.Name,
		&account.AccountNumber,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) Create(ctx context.Context, account models.Account) (*models.Account, error) {
	query := `
		INSERT INTO accounts (user_id, name, account_number, balance, currency)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	accountNumber := generateAccountNumber() // Implement this function to generate unique account numbers

	if account.Currency == "" {
		account.Currency = "NPR"
	}

	err := r.db.QueryRow(
		ctx,
		query,
		account.UserID,
		account.Name,
		accountNumber,
		account.Balance,
		account.Currency,
	).Scan(
		&account.ID,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	account.AccountNumber = accountNumber
	return &account, nil
}

func (r *AccountRepository) Update(ctx context.Context, account models.Account) error {
	query := `
		UPDATE accounts
		SET name = $1,
		    balance = $2,
		    currency = $3,
		    updated_at = NOW()
		WHERE id = $4
	`

	result, err := r.db.Exec(
		ctx,
		query,
		account.Name,
		account.Balance,
		account.Currency,
		account.ID,
	)

	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("account not found")
	}
	return nil
}
func generateAccountNumber() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("ACC%d%d", time.Now().Unix(), rand.Intn(10000))
}

func (r *AccountRepository) TransferTx(ctx context.Context, fromAccountID, toAccountID int64, amount float64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(
		ctx,
		`
			SELECT id, balance
			FROM accounts
			WHERE id = $1 OR id = $2
			ORDER BY id FOR UPDATE
		`,
		fromAccountID,
		toAccountID,
	)
	if err != nil {
		return err
	}

	defer rows.Close()

	balances := make(map[int64]float64)

	for rows.Next() {
		var id int64
		var balance float64

		err := rows.Scan(&id, &balance)
		if err != nil {
			return err
		}

		balances[id] = balance
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if _, ok := balances[fromAccountID]; !ok {
		return errors.New("from account not found")
	}

	if _, ok := balances[toAccountID]; !ok {
		return errors.New("to account not found")
	}

	if balances[fromAccountID] < amount {
		return errors.New("insufficient balance")
	}

	debitResult, err := tx.Exec(
		ctx,
		`
			UPDATE accounts
			SET balance = balance - $1, updated_at = NOW()
			WHERE id = $2
		`,
		amount,
		fromAccountID,
	)

	if err != nil {
		return err
	}

	if debitResult.RowsAffected() == 0 {
		return errors.New("from account not found")
	}

	creditResult, err := tx.Exec(
		ctx,
		`
			UPDATE accounts
			SET balance = balance + $1, updated_at = NOW()
			WHERE id = $2
		`,
		amount,
		toAccountID,
	)
	if err != nil {
		return err
	}
	if creditResult.RowsAffected() == 0 {
		return errors.New("to account not found")
	}

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO transactions (
			type,
			from_account_id,
			to_account_id,
			amount,
			status,
			reference_number
			)
			VALUES ($1, $2, $3, $4, $5, $6)
		`,
		"transfer",
		fromAccountID,
		toAccountID,
		amount,
		"success",
		generateReferenceNumber(),
	)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
