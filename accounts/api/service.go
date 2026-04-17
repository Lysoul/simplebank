package accounts

import "context"

type Service interface {
	GetAccount(ctx context.Context, id int64) (*Account, error)
	CreateAccount(ctx context.Context, owner string, balance int64, currency string) (*Account, error)
	ListAccounts(ctx context.Context) ([]*Account, error)
	DeleteAccount(ctx context.Context, id int64) error
}
