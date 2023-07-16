package simplebank

import (
	"context"
	"database/sql"
	"testing"
	"fmt"
	"github.com/stretchr/testify/require"
)


func TestTransferMoneyTx(t *testing.T){
	account_from := createAccounts(t)
	account_to := createAccounts(t)

	fmt.Printf(">>before transaction: from_account=%v, to_account=%v\n", account_from.Balance, account_to.Balance )

	store := NewStore(testDB)

	n := 5
	amount := sql.NullInt64{Int64:10, Valid:true}
	var results = make(chan TransferTxResult)
	var errs = make(chan error)

	for i:=0; i<n; i++{
		txName := fmt.Sprintf("tx%v", i+1)
		ctx := context.WithValue(context.Background(), txKey, txName)
		go func (){
			result, err := store.TransferMoneyTx(ctx, TransferMoneyTxParams{
				FromAccountID: account_from.ID,
				ToAccountID: account_to.ID,
				Amount: amount,
			})
			errs <- err
			results <- result
		}()
	}

	for i:=0; i<n; i++{
		err := <- errs
		result := <- results
		require.NoError(t, err )
		require.NotEmpty(t, result)

		transfer := result.Transfer

		require.NotEmpty(t, transfer)
		require.Equal(t, account_from.ID, transfer.FromAccountID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err )

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		require.Equal(t, -amount.Int64, fromEntry.Amount.Int64)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		require.Equal(t, amount.Int64, toEntry.Amount.Int64)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err )

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account_from.ID)
		require.NotZero(t, fromAccount.Balance)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account_to.ID)
		require.NotZero(t, toAccount.Balance)

		fmt.Printf(">>after tx%v: from_account=%v, to_account=%v\n", i+1, fromAccount.Balance, toAccount.Balance )
		// check Account balance
		diff1 := account_from.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account_to.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1>0)
		require.True(t, diff1%amount.Int64==0)

		k := diff1/amount.Int64

		require.True(t, k>=1 && k<=int64(n))
	}

	updatedFromAccount, err := testQueries.GetAccount(context.Background(), account_from.ID)

	require.NoError(t, err)
	
	updatedToAccount, err := testQueries.GetAccount(context.Background(), account_to.ID)
	require.NoError(t, err)

	fmt.Printf(">>after transaction: from_account=%v, to_account=%v\n", updatedFromAccount.Balance, updatedToAccount.Balance )
	require.Equal(t, updatedFromAccount.Balance, account_from.Balance-(int64(n)* amount.Int64) )
	require.Equal(t, updatedToAccount.Balance, account_to.Balance + (int64(n)* amount.Int64))
}

func TestTransferMoneyTxDeadlock(t *testing.T){
	account_from := createAccounts(t)
	account_to := createAccounts(t)

	fmt.Printf(">>before transaction: from_account=%v, to_account=%v\n", account_from.Balance, account_to.Balance )

	store := NewStore(testDB)

	n := 10
	amount := sql.NullInt64{Int64:10, Valid:true}
	var errs = make(chan error)

	for i:=0; i<n; i++{
		txName := fmt.Sprintf("tx%v", i+1)
		ctx := context.WithValue(context.Background(), txKey, txName)
		fromAccount := account_from.ID
		toAccount := account_to.ID

		if(i%2==1){
			fromAccount = account_to.ID
			toAccount = account_from.ID
		}

		go func (){
			_, err := store.TransferMoneyTx(ctx, TransferMoneyTxParams{
				FromAccountID: fromAccount,
				ToAccountID: toAccount,
				Amount: amount,
			})
			errs <- err
		}()
	}

	for i:=0; i<n; i++{
		err := <- errs
	
		require.NoError(t, err )
	}

	updatedFromAccount, err := testQueries.GetAccount(context.Background(), account_from.ID)

	require.NoError(t, err)
	
	updatedToAccount, err := testQueries.GetAccount(context.Background(), account_to.ID)
	require.NoError(t, err)

	fmt.Printf(">>after transaction: from_account=%v, to_account=%v\n", updatedFromAccount.Balance, updatedToAccount.Balance )
	require.Equal(t, updatedFromAccount.Balance, account_from.Balance)
	require.Equal(t, updatedToAccount.Balance, account_to.Balance)
}