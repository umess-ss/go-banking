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

func (s *AccountService) GetAccounts(ctx context.Context, userID int64) ([]models.Account, error) {
	return s.repo.FindAllByUserID(ctx, userID)
}

func (s *AccountService) GetAccountByID(ctx context.Context, accountID int64, userID int64) (*models.Account, error) {
	return s.repo.FindByIDAndUserID(ctx, accountID, userID)
}

func (s *AccountService) CreateAccount(ctx context.Context, userID int64, account models.Account) (*models.Account, error) {
	if account.Name == "" {
		return nil, errors.New("account name is required")
	}

	if account.AccountType == "" {
		account.AccountType = "savings"
	}

	if account.AccountType != "savings" && account.AccountType != "checking" && account.AccountType != "current" {
		return nil, errors.New("invalid account type")
	}

	if account.Balance < 0 {
		return nil, errors.New("initial balance cannot be negative")
	}

	account.UserID = &userID

	return s.repo.Create(ctx, account)
}

func (s *AccountService) Deposit(ctx context.Context, userID int64, accountID int64, amount float64) (*models.Account, error) {
	if amount <= 0 {
		return nil, errors.New("deposit amount must be positive")
	}

	account, err := s.repo.FindByIDAndUserID(ctx, accountID, userID)
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

func (s *AccountService) Withdraw(ctx context.Context, userID int64, accountID int64, amount float64) (*models.Account, error) {
	if amount <= 0 {
		return nil, errors.New("withdraw amount must be greater than zero")
	}

	account, err := s.repo.FindByIDAndUserID(ctx, accountID, userID)
	if err != nil {
		return nil, err
	}

	if account.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	account.Balance -= amount

	err = s.repo.Update(ctx, *account)
	if err != nil {
		return nil, err
	}

	fromAccountID := account.ID

	_, err = s.transactionRepo.Create(ctx, models.Transaction{
		Type:          "withdraw",
		FromAccountID: &fromAccountID,
		Amount:        amount,
		Status:        "success",
	})
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *AccountService) Transfer(ctx context.Context, userID int64, fromAccountID int64, toAccountID int64, amount float64) error {
	if fromAccountID == toAccountID {
		return errors.New("cannot transfer to the same account")
	}

	if amount <= 0 {
		return errors.New("transfer amount must be greater than zero")
	}

	_, err := s.repo.FindByIDAndUserID(ctx, fromAccountID, userID)
	if err != nil {
		return errors.New("from account not found")
	}

	return s.repo.TransferTx(ctx, fromAccountID, toAccountID, amount)
}
