package transaction

import (
	"context"
)

type TransactionService struct {
	repo *TransactionRepository
}

func NewTransactionService(repo *TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) GetTransactions(ctx context.Context, userID int64) ([]Transaction, error) {
	return s.repo.FindAllByUserID(ctx, userID)
}

func (s *TransactionService) GetTransactionsByAccountID(ctx context.Context, accountID int64, userID int64) ([]Transaction, error) {
	return s.repo.FindByAccountIDAndUserID(ctx, accountID, userID)
}
