package accounts

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model `json:"-"`
	Owner      string `json:"owner"`
	Balance    int64  `json:"balance"`
	Currency   string `json:"currency"`

	AccountID int64 `json:"accountId,omitempty" gorm:"-:all"` // for response only
}

type Entry struct {
	gorm.Model `json:"-"`
	AccountID  int64 `json:"accountId"`
	// can be negative or positive
	Amount int64 `json:"amount"`
}

type Transfer struct {
	gorm.Model    `json:"-"`
	FromAccountID int64 `json:"fromAccountId"`
	ToAccountID   int64 `json:"toAccountId"`
	// must be positive
	Amount int64 `json:"amount"`
}
