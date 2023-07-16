package simplebank

import (
	"context"
	"database/sql"
	"testing"

	"github.com/AlhajiTaibu/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createTransfers(t *testing.T) Transfer {
	account_from := createAccounts(t)
	account_to := createAccounts(t)

	args := CreateTransferParams{
		FromAccountID: account_from.ID,
		ToAccountID: account_to.ID,
		Amount: sql.NullInt64{Int64:util.RandomMoney(), Valid: true},
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), args)

	require.NoError(t, err)
	require.Equal(t, transfer.Amount, args.Amount)
	require.Equal(t, transfer.FromAccountID, args.FromAccountID)
	require.Equal(t, transfer.ToAccountID, args.ToAccountID)
	require.NotZero(t, transfer.CreatedAt)
	require.NotZero(t, transfer.ID)
	return transfer
}

func TestCreateTransfer(t *testing.T){
	createTransfers(t)
}

func TestGetTransfer(t *testing.T){
	transfer1 := createTransfers(t)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer1.ID, transfer2.ID)

}

func TestDeleteTransfer(t *testing.T){
	transfer := createTransfers(t)

	err := testQueries.DeleteTransfer(context.Background(), transfer.ID)

	require.NoError(t, err)
}

func TestListTransfer(t *testing.T){
	
	for i:=0; i<10; i++ {
		createTransfers(t)
	}

	args := ListTransfersParams{
		Limit: 10,
		Offset: 0,
	}
	transfers, err := testQueries.ListTransfers(context.Background(), args)

	require.NoError(t, err)

	for _, transfer := range transfers{
		require.NotEmpty(t, transfer)
	}
}

func TestUpdateTransfer(t *testing.T){
	transfer1 := createTransfers(t)

	args := UpdateTransferParams{
		ID: transfer1.ID,
		FromAccountID: transfer1.FromAccountID,
		ToAccountID: transfer1.ToAccountID,
		Amount: sql.NullInt64{Int64:util.RandomMoney(), Valid:true},
	}

	transfer2, err := testQueries.UpdateTransfer(context.Background(), args)

	require.NoError(t, err)
	require.Equal(t, transfer2.Amount, args.Amount)
	require.NotEqual(t, transfer2.Amount, transfer1.Amount)
}