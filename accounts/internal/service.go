package internal

import (
	accounts "assignments/simplebank/accounts/api"
	"assignments/simplebank/adapters/monitoring"
	"context"

	"go.uber.org/zap"
)

type Service struct {
	store accounts.Store
}

func NewService(store accounts.Store) accounts.Service {
	return &Service{
		store: store,
	}
}

// CreateAccount implements accounts.Service.
func (s *Service) CreateAccount(ctx context.Context, owner string, balance int64, currency string) (*accounts.Account, error) {
	if balance < 0 {
		monitoring.Logger().Warn("Negative balance is not allowed",
			zap.String("owner", owner),
			zap.Int64("balance", balance),
		)
	}

	if currency == "" {
		monitoring.Logger().Warn("Currency is required",
			zap.String("owner", owner),
		)
		return nil, accounts.ErrCurrencyRequired
	}

	// For simplicity, we only support USD and THB in this example.
	if currency != "USD" && currency != "THB" {
		monitoring.Logger().Warn("Unsupported currency",
			zap.String("owner", owner),
			zap.String("currency", currency),
		)
		return nil, accounts.ErrCurrencyNotSupported
	}

	return s.store.CreateAccount(ctx, owner, balance, currency)
}

// DeleteAccount implements accounts.Service.
func (s *Service) DeleteAccount(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// GetAccount implements accounts.Service.
func (s *Service) GetAccount(ctx context.Context, id int64) (*accounts.Account, error) {
	panic("unimplemented")
}

// ListAccounts implements accounts.Service.
func (s *Service) ListAccounts(ctx context.Context) ([]*accounts.Account, error) {
	panic("unimplemented")
}
