package accounts

import "assignments/simplebank/shared"

var (
	ErrNegativeBalance      = shared.ConstError("negative_balance")
	ErrCurrencyRequired     = shared.ConstError("currency_required")
	ErrCurrencyNotSupported = shared.ConstError("currency_not_supported")
)
