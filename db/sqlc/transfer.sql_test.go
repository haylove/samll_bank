//  *@createTime    2022/3/20 2:20
//  *@author        hay&object
//  *@version       v1.0.0

package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"small_bank/util"
	"testing"
)

func createRandomTransfer(t *testing.T, account1, account2 Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.CreatedAt)
	return transfer
}

func TestQueries_CreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	createRandomTransfer(t, account1, account2)
}

func TestQueries_GetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	transfer := createRandomTransfer(t, account1, account2)

	t2, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, t2)
	require.Equal(t, transfer.FromAccountID, t2.FromAccountID)
	require.Equal(t, transfer.ToAccountID, t2.ToAccountID)
	require.Equal(t, transfer.Amount, t2.Amount)
}

func TestListTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, account1, account2)
	}
	arg1 := ListOutTransferParams{
		FromAccountID: account1.ID,
		Limit:         5,
		Offset:        5,
	}
	transfer1, err := testQueries.ListOutTransfer(context.Background(), arg1)
	require.NoError(t, err)
	require.Len(t, transfer1, 5)
	for _, transfer := range transfer1 {
		require.NotEmpty(t, transfer)
	}

	arg2 := ListInTransferParams{
		ToAccountID: account2.ID,
		Limit:       5,
		Offset:      5,
	}
	transfer2, err := testQueries.ListInTransfer(context.Background(), arg2)
	require.NoError(t, err)
	require.Len(t, transfer2, 5)
	for _, transfer := range transfer2 {
		require.NotEmpty(t, transfer)
	}
	arg3 := ListUnitedTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Limit:         5,
		Offset:        5,
	}
	transfer3, err := testQueries.ListUnitedTransfer(context.Background(), arg3)
	require.NoError(t, err)
	require.Len(t, transfer3, 5)
	for _, transfer := range transfer3 {
		require.NotEmpty(t, transfer)
	}

}
