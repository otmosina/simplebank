package db

import (
	"context"
	"testing"
	"time"

	"github.com/otmosina/simplebank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomTransfer(t *testing.T, account1, account2 Account) Transfer {
	// time.Sleep(2 * time.Second)
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}
	transfer1, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer1)

	require.Equal(t, account1.ID, transfer1.FromAccountID, "should be eq")
	require.Equal(t, account2.ID, transfer1.ToAccountID, "should be eq")
	require.Equal(t, transfer1.Amount, arg.Amount)
	require.NotZero(t, transfer1.ID)
	require.WithinDuration(t, account1.CreatedAt, transfer1.CreatedAt, time.Second)
	require.WithinDuration(t, account2.CreatedAt, transfer1.CreatedAt, time.Second)

	return transfer1

}

func TestCreateTransfer(t *testing.T) {
	CreateRandomTransfer(t, CreateRandomAccount(t), CreateRandomAccount(t))
}

func TestGetTransfer(t *testing.T) {
	transfer1 := CreateRandomTransfer(t, CreateRandomAccount(t), CreateRandomAccount(t))
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)

	require.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt)
}

func TestUpdateTransfer(t *testing.T) {
	transfer1 := CreateRandomTransfer(t, CreateRandomAccount(t), CreateRandomAccount(t))
	arg := UpdateTransferParams{
		ID:     transfer1.ID,
		Amount: util.RandomMoney(),
	}
	testQueries.UpdateTransfer(context.Background(), arg)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer2.Amount, arg.Amount)
}

func TestDeletetransfer(t *testing.T) {
	transfer1 := CreateRandomEntity(t, CreateRandomAccount(t))
	testQueries.DeleteTransfer(context.Background(), transfer1.ID)
	_, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.Error(t, err)
}
func TestListTransfers(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	for i := 0; i < 10; i++ {
		CreateRandomTransfer(t, account1, account2)
	}

	arg := ListTransferParams{
		Limit:  5,
		Offset: 5,
	}
	entries, err := testQueries.ListTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)
}
