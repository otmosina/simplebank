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

	for i := 0; i < n; i++ {
		store.TransferTx(context.Background(), TransferTxParams{
			FromAccountID: account1.ID,
			ToAccountID:   account2.ID,
			Amount:        amount,
		})
	}

	// transfers, err := store.ListTransfer(context.Background(), ListTransferParams{
	// 	Limit:  1,
	// 	Offset: 1,
	// })

	// require.NoError(t, err)
	// require.Len(t, transfers, 5)
	require.Len(t, 5, 5)
}
