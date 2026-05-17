package services

import (
	"context"
	"go-banking/internal/models"
	"go-banking/internal/repository"
)

type TransactionService struct {
	repo *repository.TransactionRepository
}

func NewTransactionService(repo *repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) GetTransactions(ctx context.Context, userID int64) ([]models.Transaction, error) {
	return s.repo.FindAllByUserID(ctx, userID)
}

func (s *TransactionService) GetTransactionsByAccountID(ctx context.Context, accountID int64, userID int64) ([]models.Transaction, error) {
	return s.repo.FindByAccountIDAndUserID(ctx, accountID, userID)
}
