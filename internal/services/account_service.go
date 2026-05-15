package services

import (
	"errors"
	"go-banking/internal/models"
	"go-banking/internal/repository"
)

type AccountService struct {
	repo *repository.AccountRepository
}

func NewAccountService(repo *repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) GetAllAccounts() []models.Account {
	return s.repo.FindAll()
}

func (s *AccountService) GetAccountByID(id int) (*models.Account, error) {
	account, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("account not found")
	}
	return account, nil
}

func (s *AccountService) CreateAccount(account models.Account) (models.Account, error) {

	if account.Name == "" {
		panic("account name is required")
	}
	if account.Balance < 0 {
		panic("account balance cannot be negative")
	}
	return s.repo.Create(account), nil
}
