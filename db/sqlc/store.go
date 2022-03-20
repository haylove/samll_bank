//  *@createTime    2022/3/20 2:41
//  *@author        hay&object
//  *@version       v1.0.0

package db

import (
	"context"
	"database/sql"
	"fmt"
)

//Store provides all functions to execute db queries and transactions
type Store interface {
	Querier
	TransferTx(context.Context, TransferTxParams) (TransferTxResult, error)
}

// SQLStore implement Store
type SQLStore struct {
	*Queries
	db *sql.DB
}

//NewStore create a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

//execTX execute a function within a db transaction
func (s *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if err1 := tx.Rollback(); err1 != nil {
			return fmt.Errorf("tx err:%v rollback err:%v", err, err1)
		}
		return err
	}
	return tx.Commit()
}

//TransferTxParams contains  input parameters of transfer transaction
type TransferTxParams struct {
	FromAccountID int64   `json:"from_account_id"`
	ToAccountID   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

// TransferTxResult is a result of transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs money transfer from one account to others.
// It creates a new transfer record , add account entries add update account's balance within a single db transaction
func (s *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = s.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = s.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = s.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		return nil
	})
	return result, err
}

func addMoney(ctx context.Context, q *Queries, accountId1 int64, amount1 float64, accountId2 int64, amount2 float64) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountId1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountId2,
		Amount: amount2,
	})
	return
}
