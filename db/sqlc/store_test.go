package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	fmt.Println(" >> before: ", account1.Balance, account2.Balance)
	n := 5
	amount := int64(10)

	errors := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errors <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
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

		//check accounts

		accountFrom := res.FromAccount
		require.NotEmpty(t, accountFrom)
		require.Equal(t, account1.ID, accountFrom.ID)

		accountTo := res.ToAccount
		require.NotEmpty(t, accountTo)
		require.Equal(t, account2.ID, accountTo.ID)

		//check account balances

		fmt.Println(" >> tx: ", accountFrom.Balance, accountTo.Balance)
		diff1 := account1.Balance - accountFrom.Balance
		diff2 := accountTo.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, 1 <= k && k <= n)

		require.NotContains(t, existed, k)
		existed[k] = true
		// fmt.Println(res.Transfer.ID)
	}

	//check final update of accounts

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, updatedAccount1.Balance, account1.Balance-int64(n)*amount)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.Equal(t, updatedAccount2.Balance, account2.Balance+int64(n)*amount)

	fmt.Println(" >> after: ", updatedAccount1.Balance, updatedAccount2.Balance)

	transfers, err := store.ListTransfer(context.Background(), ListTransferParams{
		Limit:  int32(n),
		Offset: 0,
	})

	require.NoError(t, err)
	require.Len(t, transfers, n)

}
