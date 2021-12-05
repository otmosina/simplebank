package db

import (
	"context"
	"testing"
	"time"

	"github.com/otmosina/simplebank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomEntity(t *testing.T, account Account) Entry {
	account1 := account
	// time.Sleep(2 * time.Second)
	arg := CreateEntryParams{
		AccountID: account1.ID,
		Amount:    util.RandomMoney(),
	}
	entry1, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry1)

	require.Equal(t, account1.ID, entry1.AccountID, "should be eq")
	require.Equal(t, entry1.Amount, arg.Amount)
	require.NotZero(t, entry1.ID)
	require.WithinDuration(t, account1.CreatedAt, entry1.CreatedAt, time.Second)

	return entry1

}

func TestCreateEntry(t *testing.T) {
	CreateRandomEntity(t, CreateRandomAccount(t))
}

func TestGetEntry(t *testing.T) {
	entry1 := CreateRandomEntity(t, CreateRandomAccount(t))
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)

	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
}

func TestUpdateEntry(t *testing.T) {
	entry1 := CreateRandomEntity(t, CreateRandomAccount(t))
	arg := UpdateEntryParams{
		ID:     entry1.ID,
		Amount: util.RandomMoney(),
	}
	testQueries.UpdateEntry(context.Background(), arg)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry2.Amount, arg.Amount)
}

func TestDeleteEntry(t *testing.T) {
	entry1 := CreateRandomEntity(t, CreateRandomAccount(t))
	testQueries.DeleteEntry(context.Background(), entry1.ID)
	_, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.Error(t, err)
}
func TestListEntries(t *testing.T) {
	account1 := CreateRandomAccount(t)
	for i := 0; i < 10; i++ {
		CreateRandomEntity(t, account1)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}
	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)
}
