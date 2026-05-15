package repository

import (
	"go-banking/internal/models"
)

type TransactionRepository struct {
	transactions []models.Transaction
	nextID       int
}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{
		transactions: []models.Transaction{},
		nextID:       1,
	}
}

func (r *TransactionRepository) Create(transaction models.Transaction) models.Transaction {
	transaction.ID = r.nextID
	r.nextID++
	r.transactions = append(r.transactions, transaction)
	return transaction
}

func (r *TransactionRepository) FindAll() []models.Transaction {
	return r.transactions
}

func (r *TransactionRepository) FindByAccountID(accountID int) []models.Transaction {
	result := []models.Transaction{}

	for _, transaction := range r.transactions {
		fromMatch := transaction.FromAccountID != nil && *transaction.FromAccountID == accountID
		toMatch := transaction.ToAccountID != nil && *transaction.ToAccountID == accountID

		if fromMatch || toMatch {
			result = append(result, transaction)
		}
	}

	return result
}
