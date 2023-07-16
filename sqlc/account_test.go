package simplebank

import (
	"github.com/AlhajiTaibu/simplebank/util"
	"context"
	"testing"
	"github.com/stretchr/testify/require"
)

func createAccounts(t *testing.T) Account{
	args := CreateAccountParams{
	Owner:    util.RandomOwner(),
	Currency: util.RandomCurrency(),
	Balance:  util.RandomMoney(),
	}

	account, err := testQueries.CreateAccount(context.Background(), args)

	require.NoError(t, err, "No errors")
	require.NotEmpty(t, account)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createAccounts(t)
}

func TestGetAccount(t *testing.T){
	account1:=createAccounts(t)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err, "No errors")
	require.NotEmpty(t, account2)
}

func TestUpdateAccount(t *testing.T){
	account1:=createAccounts(t)

	args := UpdateAccountParams{
		ID: account1.ID,
		Owner: util.RandomOwner(),
		Balance: util.RandomMoney(),
	}

	account2, err:= testQueries.UpdateAccount(context.Background(), args)

	require.NoError(t, err, "No errors")
	require.NotEmpty(t, account2)
	require.Equal(t, args.Balance, account2.Balance)
	require.Equal(t, args.Owner, account2.Owner)
	require.Equal(t, args.Currency, account2.Currency)
}

func TestListAccounts(t *testing.T){
	for i:=0; i<10; i++{
		createAccounts(t)
	}

	accounts, err := testQueries.ListAccounts(context.Background())

	require.NoError(t, err, "No errors")

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestDeleteAccount(t *testing.T){
	account := createAccounts(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)

	require.NoError(t, err)
}
