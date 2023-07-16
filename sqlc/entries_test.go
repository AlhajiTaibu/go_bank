package simplebank

import (
	"context"
	"database/sql"
	"testing"

	"github.com/AlhajiTaibu/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createEntries(t *testing.T) Entry{
	account := createAccounts(t)
	args := CreateEntryParams{
		AccountID: account.ID,
		Amount: sql.NullInt64{Int64:util.RandomMoney(), Valid: true},
	}

	entry, err := testQueries.CreateEntry(context.Background(), args)

	require.NoError(t, err, "No errors")
	require.NotEmpty(t, account)
	require.Equal(t, args.Amount, entry.Amount)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T){
	createEntries(t)
}

func TestGetEntry(t *testing.T){
	entry1 := createEntries(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err, "No errors")
	require.NotEmpty(t, entry2)
}

func TestUpdateEntry(t *testing.T){
	entry1 := createEntries(t)
	args := UpdateEntryParams{
		ID: entry1.ID,
		AccountID: entry1.AccountID,
		Amount: sql.NullInt64{Int64:util.RandomMoney(), Valid: true},
	}

	entry2, err := testQueries.UpdateEntry(context.Background(), args)

	require.NoError(t, err, "No errors")
	require.Equal(t, args.Amount, entry2.Amount)
	require.NotZero(t, entry2.ID)
}

func TestListEntires(t *testing.T){
	for i:=0; i<10; i++{
		createEntries(t)
	}

	entries, err := testQueries.ListEntries(context.Background())

	require.NoError(t, err, "No errors")

	for _, entry := range entries{
		require.NotEmpty(t, entry)
	}
}

func TestDelete(t *testing.T){
	entry := createEntries(t)

	err := testQueries.DeleteEntry(context.Background(), entry.ID)
	
	require.NoError(t, err)
}