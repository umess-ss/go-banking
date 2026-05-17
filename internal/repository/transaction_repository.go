package repository

import (
	"context"
	"fmt"
	"go-banking/internal/models"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) Create(ctx context.Context, transaction models.Transaction) (models.Transaction, error) {
	query := `
		INSERT INTO transactions (
		type, from_account_id, to_account_id, amount, status, reference_number)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, type, from_account_id, to_account_id, amount, status, reference_number, created_at
	`

	if transaction.Status == "" {
		transaction.Status = "success"
	}

	if transaction.ReferenceNumber == "" {
		transaction.ReferenceNumber = generateReferenceNumber()
	}

	err := r.db.QueryRow(
		ctx,
		query,
		transaction.Type,
		transaction.FromAccountID,
		transaction.ToAccountID,
		transaction.Amount,
		transaction.Status,
		transaction.ReferenceNumber,
	).Scan(
		&transaction.ID,
		&transaction.Type,
		&transaction.FromAccountID,
		&transaction.ToAccountID,
		&transaction.Amount,
		&transaction.Status,
		&transaction.ReferenceNumber,
		&transaction.CreatedAt,
	)
	if err != nil {
		return models.Transaction{}, err
	}
	return transaction, nil
}

func (r *TransactionRepository) FindAllByUserID(ctx context.Context, userID int64) ([]models.Transaction, error) {
	query := `
		SELECT DISTINCT t.id, t.type, t.from_account_id, t.to_account_id, t.amount, t.status, t.reference_number, t.created_at
		FROM transactions t
			LEFT JOIN accounts from_acc ON t.from_account_id = from_acc.id
			LEFT JOIN accounts to_acc ON t.to_account_id = to_acc.id
		WHERE from_acc.user_id = $1 OR to_acc.user_id = $1
		ORDER BY t.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []models.Transaction{}

	for rows.Next() {
		var transaction models.Transaction

		err := rows.Scan(
			&transaction.ID,
			&transaction.Type,
			&transaction.FromAccountID,
			&transaction.ToAccountID,
			&transaction.Amount,
			&transaction.Status,
			&transaction.ReferenceNumber,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionRepository) FindByAccountIDAndUserID(ctx context.Context, accountID int64, userID int64) ([]models.Transaction, error) {
	query := `
		SELECT t.id, t.type, t.from_account_id, t.to_account_id, t.amount, t.status, t.reference_number, t.created_at
		FROM transactions t
		INNER JOIN accounts a ON a.id = $1
		WHERE a.user_id = $2
		  AND (t.from_account_id = $1 OR t.to_account_id = $1)
		ORDER BY t.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, accountID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []models.Transaction{}

	for rows.Next() {
		var transaction models.Transaction

		err := rows.Scan(
			&transaction.ID,
			&transaction.Type,
			&transaction.FromAccountID,
			&transaction.ToAccountID,
			&transaction.Amount,
			&transaction.Status,
			&transaction.ReferenceNumber,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func generateReferenceNumber() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("TXN%d%d", time.Now().UnixNano(), rand.Intn(10000))
}
