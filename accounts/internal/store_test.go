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

func TestGetAccount(t *testing.T) {

	account := &api.Account{
		Owner:    "Alice",
		Balance:  2000,
		Currency: "USD",
	}

	result, err := store.CreateAccount(context.Background(), account.Owner, account.Balance, account.Currency)
	require.NoError(t, err)

	result, err = store.GetAccount(context.Background(), int64(result.ID))
	require.NoError(t, err)
	require.Equal(t, account.Owner, result.Owner)
	require.Equal(t, account.Balance, result.Balance)
	require.Equal(t, account.Currency, result.Currency)
}

func TestGetAccountNotFound(t *testing.T) {
	_, err := store.GetAccount(context.Background(), 99999)
	require.ErrorIs(t, err, api.ErrAccountNotFound)
}

func TestListAccounts(t *testing.T) {
	account1 := &api.Account{
		Owner:    "Alice",
		Balance:  2000,
		Currency: "USD",
	}
	account2 := &api.Account{
		Owner:    "Bob",
		Balance:  1000,
		Currency: "USD",
	}

	_, err := store.CreateAccount(context.Background(), account1.Owner, account1.Balance, account1.Currency)
	require.NoError(t, err)
	_, err = store.CreateAccount(context.Background(), account2.Owner, account2.Balance, account2.Currency)
	require.NoError(t, err)

	results, err := store.ListAccounts(context.Background())
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(results), 2)
}

func TestDeleteAccount(t *testing.T) {
	account := &api.Account{
		Owner:    "Alice",
		Balance:  2000,
		Currency: "USD",
	}

	result, err := store.CreateAccount(context.Background(), account.Owner, account.Balance, account.Currency)
	require.NoError(t, err)

	err = store.DeleteAccount(context.Background(), int64(result.ID))
	require.NoError(t, err)

	_, err = store.GetAccount(context.Background(), int64(result.ID))
	require.ErrorIs(t, err, api.ErrAccountNotFound)
}

func TestTransferTx(t *testing.T) {

	account1 := &api.Account{
		Owner:    "Alice",
		Balance:  2000,
		Currency: "USD",
	}

	account2 := &api.Account{
		Owner:    "Bob",
		Balance:  1000,
		Currency: "USD",
	}

	result1, err := store.CreateAccount(context.Background(), account1.Owner, account1.Balance, account1.Currency)
	require.NoError(t, err)

	result2, err := store.CreateAccount(context.Background(), account2.Owner, account2.Balance, account2.Currency)
	require.NoError(t, err)

	transfer, err := store.TransferTx(context.Background(), int64(result1.ID), int64(result2.ID), 500)
	require.NoError(t, err)
	require.NotNil(t, transfer)
	require.Equal(t, int64(result1.ID), transfer.FromAccountID)
	require.Equal(t, int64(result2.ID), transfer.ToAccountID)
	require.Equal(t, int64(500), transfer.Amount)
}
