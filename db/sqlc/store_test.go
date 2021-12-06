package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	n := 5
	amount := int64(10)

	errors := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errors <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)

		res := <-results
		require.NotEmpty(t, res)

		// check transfer
		transfer := res.Transfer

		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, errGetTransfer := testQueries.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, errGetTransfer)

		//check Entry from

		fromEntry := res.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, errGetEntry := testQueries.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, errGetEntry)

		//check EntryTo

		toEntry := res.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, errGetEntry2 := testQueries.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, errGetEntry2)

		// fmt.Println(res.Transfer.ID)
	}

	transfers, err := store.ListTransfer(context.Background(), ListTransferParams{
		Limit:  5,
		Offset: 0,
	})

	require.NoError(t, err)
	require.Len(t, transfers, 5)

}
