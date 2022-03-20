//  *@createTime    2022/3/20 3:21
//  *@author        hay&object
//  *@version       v1.0.0

package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_TransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	errChan := make(chan error)
	resultChan := make(chan TransferTxResult)

	for i := 0; i < 2; i++ {
		go func() {
			txResult, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        10,
			})
			errChan <- err
			resultChan <- txResult
		}()
	}

	for i := 0; i < 2; i++ {
		err := <-errChan
		require.NoError(t, err)

		result := <-resultChan
		require.NotEmpty(t, result)
		fromEntry := result.FromEntry
		require.Equal(t, fromEntry.Amount, float64(-10))
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		toEntry := result.ToEntry
		require.Equal(t, toEntry.Amount, float64(10))
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.NotZero(t, toEntry.CreatedAt)

		fromAccount := result.FromAccount
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.Equal(t, toAccount.ID, account2.ID)
	}
}
