//  *@createTime    2022/3/20 1:53
//  *@author        hay&object
//  *@version       v1.0.0

package db

import (
	"context"
	"testing"
	"time"

	"github.com/haylove/small_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, account.ID, entry.AccountID)
	return entry
}

func TestQueries_CreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestQueries_GetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1 := createRandomEntry(t, account)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestQueries_UpdateEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1 := createRandomEntry(t, account)
	arg := UpdateEntryParams{
		ID:     entry1.ID,
		Amount: util.RandMoney(),
	}
	entry2, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, arg.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestQueries_ListEntriesByAccountID(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}
	arg := ListEntriesByAccountIDParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}
	entries, err := testQueries.ListEntriesByAccountID(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
