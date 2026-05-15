package services

import (
	"errors"
	"go-banking/internal/models"
	"go-banking/internal/repository"
)

// service is business logic layer that interacts with the repository to perform operations
// on accounts. It provides methods for retrieving and creating accounts, and can include
// additional business rules or validations as needed.
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

func (s *AccountService) Deposit(accountID int, amount float64) (*models.Account, error) {
	if amount <= 0 {
		return nil, errors.New("deposit amount must be positive")
	}
	account, err := s.GetAccountByID(accountID)
	if err != nil {
		return nil, err
	}
	account.Balance += amount
	err = s.repo.Update(*account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *AccountService) Withdraw(accountID int, amount float64) (*models.Account, error) {
	if amount <= 0 {
		return nil, errors.New("withdrawal amount must be positive")
	}
	account, err := s.GetAccountByID(accountID)
	if err != nil {
		return nil, err
	}
	if account.Balance < amount {
		return nil, errors.New("insufficient funds")
	}
	account.Balance -= amount
	err = s.repo.Update(*account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *AccountService) Transfer(fromAccountID, toAccountID int, amount float64) error {
	if amount <= 0 {
		return errors.New("transfer amount must be positive")
	}
	fromAccount, err := s.GetAccountByID(fromAccountID)
	if err != nil {
		return err
	}
	toAccount, err := s.GetAccountByID(toAccountID)
	if err != nil {
		return err
	}
	if fromAccount.Balance < amount {
		return errors.New("insufficient funds")
	}
	fromAccount.Balance -= amount
	toAccount.Balance += amount
	err = s.repo.Update(*fromAccount)
	if err != nil {
		return err
	}
	err = s.repo.Update(*toAccount)
	if err != nil {
		return err
	}
	return nil
}
