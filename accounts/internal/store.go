package internal

import (
	accounts "assignments/simplebank/accounts/api"
	"assignments/simplebank/adapters/monitoring"
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewStore(db *gorm.DB) accounts.Store {
	store := &storeRDB{
		db,
	}
	return store
}

type storeRDB struct {
	db *gorm.DB
}

// CreateAccount implements accounts.Store.
func (s *storeRDB) CreateAccount(ctx context.Context, owner string, balance int64, currency string) (*accounts.Account, error) {
	account := &accounts.Account{
		Owner:    owner,
		Balance:  balance,
		Currency: currency,
	}

	res := s.db.Create(account)
	if res.Error != nil {
		monitoring.Logger().Error("failed to create account", zap.Error(res.Error))
		return nil, res.Error
	}

	return account, nil
}

// DeleteAccount implements accounts.Store.
func (s *storeRDB) DeleteAccount(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// GetAccount implements accounts.Store.
func (s *storeRDB) GetAccount(ctx context.Context, id int64) (*accounts.Account, error) {
	panic("unimplemented")
}

// ListAccounts implements accounts.Store.
func (s *storeRDB) ListAccounts(ctx context.Context) ([]*accounts.Account, error) {
	panic("unimplemented")
}
