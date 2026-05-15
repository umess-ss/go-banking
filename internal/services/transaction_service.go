package services

import (
	"go-banking/internal/models"
	"go-banking/internal/repository"
)

type TransactionService struct {
	repo *repository.TransactionRepository
}

func NewTransactionService(repo *repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) GetTransactions() []models.Transaction {
	return s.repo.FindAll()
}

func (s *TransactionService) GetTransactionsByAccountID(accountID int) []models.Transaction {
	return s.repo.FindByAccountID(accountID)
}
