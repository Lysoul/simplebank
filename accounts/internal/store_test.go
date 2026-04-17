package internal_test

import (
	api "assignments/simplebank/accounts/api"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	account := &api.Account{
		Owner:    "Bob",
		Balance:  1000,
		Currency: "USD",
	}

	result, err := store.CreateAccount(context.Background(), account.Owner, account.Balance, account.Currency)

	require.NoError(t, err)
	require.Equal(t, account.Owner, result.Owner)
	require.Equal(t, account.Balance, result.Balance)
	require.Equal(t, account.Currency, result.Currency)
}
