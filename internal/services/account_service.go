package services

import (
	"context"
	"errors"
	"go-banking/internal/models"
	"go-banking/internal/repository"
	"time"
)

// service is business logic layer that interacts with the repository to perform operations
// on accounts. It provides methods for retrieving and creating accounts, and can include
// additional business rules or validations as needed.
type AccountService struct {
	repo            *repository.AccountRepository
	transactionRepo *repository.TransactionRepository
}

func NewAccountService(repo *repository.AccountRepository, transactionRepo *repository.TransactionRepository) *AccountService {
	return &AccountService{
		repo:            repo,
		transactionRepo: transactionRepo,
	}
}

func (s *AccountService) GetAllAccounts(ctx context.Context) ([]models.Account, error) {
	return s.repo.FindAll(ctx)
}

func (s *AccountService) GetAccountByID(ctx context.Context, id int64) (*models.Account, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *AccountService) CreateAccount(ctx context.Context, account models.Account) (*models.Account, error) {

	if account.Name == "" {
		return nil, errors.New("account name is required")
	}
	if account.Balance < 0 {
		return nil, errors.New("account balance cannot be negative")
	}
	return s.repo.Create(ctx, account)
}

func (s *AccountService) Deposit(ctx context.Context, accountID int64, amount float64) (*models.Account, error) {
	if amount <= 0 {
		return nil, errors.New("deposit amount must be positive")
	}

	account, err := s.repo.FindByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	account.Balance += amount
	err = s.repo.Update(ctx, *account)
	if err != nil {
		return nil, err
	}

	_, err = s.transactionRepo.Create(ctx, models.Transaction{
		Type:        "deposit",
		ToAccountID: &accountID,
		Amount:      amount,
		CreatedAt:   time.Now(),
	})
	return account, nil
}

func (s *AccountService) Withdraw(ctx context.Context, accountID int64, amount float64) (*models.Account, error) {
	if amount <= 0 {
		return nil, errors.New("withdrawal amount must be positive")
	}
	account, err := s.repo.FindByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if account.Balance < amount {
		return nil, errors.New("insufficient funds")
	}
	account.Balance -= amount
	err = s.repo.Update(ctx, *account)
	if err != nil {
		return nil, err
	}

	_, err = s.transactionRepo.Create(ctx, models.Transaction{
		Type:          "withdrawal",
		FromAccountID: &accountID,
		Amount:        amount,
		CreatedAt:     time.Now(),
	})
	return account, nil
}

func (s *AccountService) Transfer(ctx context.Context, fromAccountID int64, toAccountID int64, amount float64) error {
	if fromAccountID == toAccountID {
		return errors.New("cannot transfer to the same account")
	}

	if amount <= 0 {
		return errors.New("transfer amount must be greater than zero")
	}

	fromAccount, err := s.repo.FindByID(ctx, fromAccountID)
	if err != nil {
		return errors.New("from account not found")
	}

	toAccount, err := s.repo.FindByID(ctx, toAccountID)
	if err != nil {
		return errors.New("to account not found")
	}

	if fromAccount.Balance < amount {
		return errors.New("insufficient balance")
	}

	fromAccount.Balance -= amount
	toAccount.Balance += amount

	err = s.repo.Update(ctx, *fromAccount)
	if err != nil {
		return err
	}

	err = s.repo.Update(ctx, *toAccount)
	if err != nil {
		return err
	}

	_, err = s.transactionRepo.Create(ctx, models.Transaction{
		Type:          "transfer",
		FromAccountID: &fromAccountID,
		ToAccountID:   &toAccountID,
		Amount:        amount,
		CreatedAt:     time.Now(),
	})

	return nil
}
