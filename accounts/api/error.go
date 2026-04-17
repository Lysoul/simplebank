package accounts

import "assignments/simplebank/shared"

var (
	ErrNegativeBalance      = shared.ConstError("negative_balance")
	ErrCurrencyRequired     = shared.ConstError("currency_required")
	ErrCurrencyNotSupported = shared.ConstError("currency_not_supported")
	ErrInvalidID            = shared.ConstError("invalid_id")
	ErrAccountNotFound      = shared.ConstError("account_not_found")

	ErrInvalidTransferAmount = shared.ConstError("invalid_transfer_amount")
	ErrSameAccountTransfer   = shared.ConstError("same_account_transfer")
	ErrInsufficientBalance   = shared.ConstError("insufficient_balance")
)
