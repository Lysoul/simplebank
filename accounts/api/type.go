package accounts

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Owner    string `json:"owner"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

type Entry struct {
	gorm.Model
	AccountID int64 `json:"account_id"`
	// can be negative or positive
	Amount int64 `json:"amount"`
}

type Transfer struct {
	gorm.Model
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	// must be positive
	Amount int64 `json:"amount"`
}
