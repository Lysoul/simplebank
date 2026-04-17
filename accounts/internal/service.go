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
		return nil, accounts.ErrNegativeBalance
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

	account, err := s.store.CreateAccount(ctx, owner, balance, currency)
	if err != nil {
		return nil, err
	}

	account.AccountID = int64(account.ID)
	return account, nil
}

// DeleteAccount implements accounts.Service.
func (s *Service) DeleteAccount(ctx context.Context, id int64) error {
	return s.store.DeleteAccount(ctx, id)
}

// GetAccount implements accounts.Service.
func (s *Service) GetAccount(ctx context.Context, id int64) (*accounts.Account, error) {
	return s.store.GetAccount(ctx, id)
}

// ListAccounts implements accounts.Service.
func (s *Service) ListAccounts(ctx context.Context) ([]*accounts.Account, error) {
	return s.store.ListAccounts(ctx)
}

// TransferTx implements accounts.Service.
func (s *Service) TransferTx(ctx context.Context, fromAccountID int64, toAccountID int64, amount int64) (*accounts.Transfer, error) {

	if amount <= 0 {
		monitoring.Logger().Warn("Transfer amount must be positive",
			zap.Int64("fromAccountID", fromAccountID),
			zap.Int64("toAccountID", toAccountID),
			zap.Int64("amount", amount),
		)
		return nil, accounts.ErrInvalidTransferAmount
	}

	if fromAccountID == toAccountID {
		monitoring.Logger().Warn("Cannot transfer to the same account",
			zap.Int64("accountID", fromAccountID),
		)
		return nil, accounts.ErrSameAccountTransfer
	}

	fromAccount, err := s.store.GetAccount(ctx, fromAccountID)
	if err != nil {
		if err == accounts.ErrAccountNotFound {
			monitoring.Logger().Warn("From account not found",
				zap.Int64("accountID", fromAccountID),
			)
			return nil, accounts.ErrAccountNotFound
		}
		return nil, err
	}

	_, err = s.store.GetAccount(ctx, toAccountID)
	if err != nil {
		if err == accounts.ErrAccountNotFound {
			monitoring.Logger().Warn("To account not found",
				zap.Int64("accountID", toAccountID),
			)
			return nil, accounts.ErrAccountNotFound
		}
		return nil, err
	}

	if fromAccount.Balance < amount {
		monitoring.Logger().Warn("Insufficient balance in from account",
			zap.Int64("accountID", fromAccountID),
			zap.Int64("balance", fromAccount.Balance),
			zap.Int64("transferAmount", amount),
		)
		return nil, accounts.ErrInsufficientBalance
	}

	transfer, err := s.store.TransferTx(ctx, fromAccountID, toAccountID, amount)
	if err != nil {
		return nil, err
	}

	return transfer, nil
}
