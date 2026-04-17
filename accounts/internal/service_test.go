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
		name     string
		owner    string
		balance  int64
		currency string
		err      error
	}{
		{
			name:     "success",
			owner:    "Bob",
			balance:  1000,
			currency: "USD",
			err:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mocks.NewMockStore(ctrl)
			mockStore.EXPECT().
				CreateAccount(gomock.Any(), tc.owner, tc.balance, tc.currency).
				Return(&accounts.Account{
					Owner:    tc.owner,
					Balance:  tc.balance,
					Currency: tc.currency,
				}, tc.err)

			service := internal.NewService(mockStore)
			account, err := service.CreateAccount(context.Background(), tc.owner, tc.balance, tc.currency)

			if tc.err != nil {
				require.Error(t, err, tc.err)
			}
			require.Equal(t, tc.owner, account.Owner)
			require.Equal(t, tc.balance, account.Balance)
			require.Equal(t, tc.currency, account.Currency)
		})
	}

}
