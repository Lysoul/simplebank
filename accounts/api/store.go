package accounts

import "context"

//go:generate mockgen -source=./store.go -destination=./mocks/store.go -package=mocks "assignments/simplebank" Store
type Store interface {
	GetAccount(ctx context.Context, id int64) (*Account, error)
	CreateAccount(ctx context.Context, owner string, balance int64, currency string) (*Account, error)
	ListAccounts(ctx context.Context) ([]*Account, error)
	DeleteAccount(ctx context.Context, id int64) error
}
