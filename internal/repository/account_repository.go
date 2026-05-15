package repository

import (
	"errors"
	"go-banking/internal/models"
)

//AccountRepository handles data storage and retrieval for accounts.
//In a real application, this would interface with a database, but for simplicity,
// we'll use an in-memory slice.

type AccountRepository struct {
	accounts []models.Account
	nextID   int
}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{
		accounts: []models.Account{},
		nextID:   1,
	}
}

func (r *AccountRepository) FindAll() []models.Account {
	return r.accounts
}

func (r *AccountRepository) FindByID(id int) (*models.Account, error) {
	for _, account := range r.accounts {
		if account.ID == id {
			return &account, nil
		}
	}
	return nil, errors.New("account not found")
}

func (r *AccountRepository) Create(account models.Account) models.Account {
	account.ID = r.nextID
	r.nextID++
	r.accounts = append(r.accounts, account)
	return account
}

func (r *AccountRepository) Update(account models.Account) error {
	for i := range r.accounts {
		if r.accounts[i].ID == account.ID {
			r.accounts[i] = account
			return nil
		}
	}
	return errors.New("account not found")
}
