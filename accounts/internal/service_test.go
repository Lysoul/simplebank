package internal_test

import (
	accounts "assignments/simplebank/accounts/api"
	"assignments/simplebank/accounts/api/mocks"
	"assignments/simplebank/accounts/internal"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateAccountService(t *testing.T) {
	testCases := []struct {
		name       string
		owner      string
		balance    int64
		currency   string
		buildStubs func(store *mocks.MockStore, owner string, balance int64, currency string)
		err        error
	}{
		{
			name:     "success",
			owner:    "Bob",
			balance:  1000,
			currency: "USD",
			buildStubs: func(store *mocks.MockStore,
				owner string,
				balance int64,
				currency string,
			) {
				store.EXPECT().
					CreateAccount(gomock.Any(), owner, balance, currency).
					Return(&accounts.Account{
						Owner:    owner,
						Balance:  balance,
						Currency: currency,
					}, nil).Times(1)
			},
			err: nil,
		},
		{
			name:       "negative balance",
			owner:      "Bob",
			balance:    -1000,
			currency:   "USD",
			buildStubs: func(store *mocks.MockStore, owner string, balance int64, currency string) {},
			err:        accounts.ErrNegativeBalance,
		},
		{
			name:       "missing currency",
			owner:      "Charlie",
			balance:    1000,
			currency:   "",
			buildStubs: func(store *mocks.MockStore, owner string, balance int64, currency string) {},
			err:        accounts.ErrCurrencyRequired,
		},
		{
			name:       "unsupported currency",
			owner:      "Alice",
			balance:    1000,
			currency:   "EUR",
			buildStubs: func(store *mocks.MockStore, owner string, balance int64, currency string) {},
			err:        accounts.ErrCurrencyNotSupported,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mocks.NewMockStore(ctrl)
			tc.buildStubs(mockStore, tc.owner, tc.balance, tc.currency)

			service := internal.NewService(mockStore)
			account, err := service.CreateAccount(context.Background(), tc.owner, tc.balance, tc.currency)

			if tc.err != nil {
				require.Error(t, err, tc.err)
				return
			}
			require.Equal(t, tc.owner, account.Owner)
			require.Equal(t, tc.balance, account.Balance)
			require.Equal(t, tc.currency, account.Currency)
		})
	}

}

func TestGetAccountService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	accountID := int64(1)
	mockStore.EXPECT().
		GetAccount(gomock.Any(), accountID).
		Return(&accounts.Account{
			AccountID: accountID,
			Owner:     "Bob",
			Balance:   1000,
			Currency:  "USD",
		}, nil).Times(1)

	service := internal.NewService(mockStore)
	result, err := service.GetAccount(context.Background(), accountID)
	require.NoError(t, err)
	require.Equal(t, accountID, result.AccountID)
	require.Equal(t, "Bob", result.Owner)
	require.Equal(t, int64(1000), result.Balance)
	require.Equal(t, "USD", result.Currency)
}

func TestTransferTxService(t *testing.T) {
	testCases := []struct {
		name          string
		fromAccountID int64
		toAccountID   int64
		amount        int64
		buildStubs    func(store *mocks.MockStore, fromAccountID int64, toAccountID int64, amount int64)
		err           error
	}{
		{
			name:          "success",
			fromAccountID: 1,
			toAccountID:   2,
			amount:        500,
			buildStubs: func(store *mocks.MockStore, fromAccountID int64, toAccountID int64, amount int64) {
				store.EXPECT().
					GetAccount(gomock.Any(), fromAccountID).
					Return(&accounts.Account{
						AccountID: fromAccountID,
						Owner:     "Alice",
						Balance:   2000,
						Currency:  "USD",
					}, nil).Times(1)
				store.EXPECT().
					GetAccount(gomock.Any(), toAccountID).
					Return(&accounts.Account{
						AccountID: toAccountID,
						Owner:     "Bob",
						Balance:   1000,
						Currency:  "USD",
					}, nil).Times(1)
				store.EXPECT().
					TransferTx(gomock.Any(), fromAccountID, toAccountID, amount).
					Return(&accounts.Transfer{
						FromAccountID: fromAccountID,
						ToAccountID:   toAccountID,
						Amount:        amount,
					}, nil).Times(1)
			},
			err: nil,
		},
		{
			name:          "invalid transfer amount",
			fromAccountID: 1,
			toAccountID:   2,
			amount:        -500,
			buildStubs:    func(store *mocks.MockStore, fromAccountID int64, toAccountID int64, amount int64) {},
			err:           accounts.ErrInvalidTransferAmount,
		},
		{
			name:          "same account transfer",
			fromAccountID: 1,
			toAccountID:   1,
			amount:        500,
			buildStubs:    func(store *mocks.MockStore, fromAccountID int64, toAccountID int64, amount int64) {},
			err:           accounts.ErrSameAccountTransfer,
		},
		{
			name:          "from account not found",
			fromAccountID: 1,
			toAccountID:   2,
			amount:        500,
			buildStubs: func(store *mocks.MockStore, fromAccountID int64, toAccountID int64, amount int64) {
				store.EXPECT().
					GetAccount(gomock.Any(), fromAccountID).
					Return(nil, accounts.ErrAccountNotFound).Times(1)
			},
			err: accounts.ErrAccountNotFound,
		},
		{
			name:          "to account not found",
			fromAccountID: 1,
			toAccountID:   2,
			amount:        500,
			buildStubs: func(store *mocks.MockStore, fromAccountID int64, toAccountID int64, amount int64) {
				store.EXPECT().
					GetAccount(gomock.Any(), fromAccountID).
					Return(&accounts.Account{
						AccountID: fromAccountID,
						Owner:     "Alice",
						Balance:   2000,
						Currency:  "USD",
					}, nil).Times(1)
				store.EXPECT().
					GetAccount(gomock.Any(), toAccountID).
					Return(nil, accounts.ErrAccountNotFound).Times(1)
			},
			err: accounts.ErrAccountNotFound,
		},
		{
			name:          "insufficient balance",
			fromAccountID: 1,
			toAccountID:   2,
			amount:        5000,
			buildStubs: func(store *mocks.MockStore, fromAccountID int64, toAccountID int64, amount int64) {
				store.EXPECT().
					GetAccount(gomock.Any(), fromAccountID).
					Return(&accounts.Account{
						AccountID: fromAccountID,
						Owner:     "Alice",
						Balance:   2000,
						Currency:  "USD",
					}, nil).Times(1)

				store.EXPECT().
					GetAccount(gomock.Any(), toAccountID).
					Return(&accounts.Account{
						AccountID: toAccountID,
						Owner:     "Bob",
						Balance:   1000,
						Currency:  "USD",
					}, nil).Times(1)
			},
			err: accounts.ErrInsufficientBalance,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mocks.NewMockStore(ctrl)
			tc.buildStubs(mockStore, tc.fromAccountID, tc.toAccountID, tc.amount)

			service := internal.NewService(mockStore)
			transfer, err := service.TransferTx(context.Background(), tc.fromAccountID, tc.toAccountID, tc.amount)

			if tc.err != nil {
				require.Error(t, err, tc.err)
				return
			}
			require.Equal(t, tc.fromAccountID, transfer.FromAccountID)
			require.Equal(t, tc.toAccountID, transfer.ToAccountID)
			require.Equal(t, tc.amount, transfer.Amount)
		})

	}
}
