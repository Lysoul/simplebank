package internal

import (
	accounts "assignments/simplebank/accounts/api"
	"assignments/simplebank/adapters/monitoring"
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// TransferTx implements accounts.Store.
func (s *storeRDB) TransferTx(ctx context.Context, fromAccountID int64, toAccountID int64, amount int64) (*accounts.Transfer, error) {
	var transfer *accounts.Transfer

	err := s.db.Transaction(func(tx *gorm.DB) error {

		firstID := fromAccountID
		secondID := toAccountID

		//lock ordering to prevent deadlock
		if fromAccountID > toAccountID {
			firstID = toAccountID
			secondID = fromAccountID
		}

		var acc1, acc2 accounts.Account

		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", firstID).
			First(&acc1).Error; err != nil {
			return err
		}

		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", secondID).
			First(&acc2).Error; err != nil {
			return err
		}

		transfer = &accounts.Transfer{
			FromAccountID: fromAccountID,
			ToAccountID:   toAccountID,
			Amount:        amount,
		}
		if err := tx.Create(transfer).Error; err != nil {
			return err
		}

		fromEntry := &accounts.Entry{
			AccountID: fromAccountID,
			Amount:    -amount,
		}
		if err := tx.Create(fromEntry).Error; err != nil {
			return err
		}

		toEntry := &accounts.Entry{
			AccountID: toAccountID,
			Amount:    amount,
		}
		if err := tx.Create(toEntry).Error; err != nil {
			return err
		}

		if err := tx.Model(&accounts.Account{}).
			Where("id = ?", fromAccountID).
			Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}

		if err := tx.Model(&accounts.Account{}).
			Where("id = ?", toAccountID).
			Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}

		return nil
	})

	return transfer, err
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
	res := s.db.Delete(&accounts.Account{}, id)
	if res.Error != nil {
		monitoring.Logger().Error("failed to delete account", zap.Int64("id", id), zap.Error(res.Error))
		return res.Error
	}
	return nil
}

// GetAccount implements accounts.Store.
func (s *storeRDB) GetAccount(ctx context.Context, id int64) (*accounts.Account, error) {
	var account accounts.Account
	res := s.db.Where("id = ?", id).First(&account)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, accounts.ErrAccountNotFound
		}

		monitoring.Logger().Error("failed to get account", zap.Int64("id", id), zap.Error(res.Error))
		return nil, res.Error
	}
	return &account, nil

}

// ListAccounts implements accounts.Store.
func (s *storeRDB) ListAccounts(ctx context.Context) ([]*accounts.Account, error) {
	var accounts []*accounts.Account
	//todo: add pagination later
	res := s.db.Find(&accounts)
	if res.Error != nil {
		monitoring.Logger().Error("failed to list accounts", zap.Error(res.Error))
		return nil, res.Error
	}
	return accounts, nil
}
