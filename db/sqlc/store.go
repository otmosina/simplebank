package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil) //&sql.TxOptions{})
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {

		txName := ctx.Value(txKey)
		var err error

		fmt.Println("create transfer", txName)
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Println("create from_entry", txName)
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println("create to_entry", txName)
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// TODO: udpate account balances

		// var account1, account2 Account
		// fmt.Println("get account1", txName)
		// account1, err = q.GetAccountForUpdate(ctx, arg.FromAccountID)
		// if err != nil {
		// 	return err
		// }

		fmt.Println("update account2", txName)
		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		// result.FromAccount, err = q.GetAccount(ctx, arg.FromAccountID)
		// if err != nil {
		// 	return err
		// }

		// fmt.Println("get account2", txName)
		// account2, err = q.GetAccountForUpdate(ctx, arg.ToAccountID)
		// if err != nil {
		// 	return err
		// }

		fmt.Println("update account2", txName)
		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		// result.ToAccount, err = q.GetAccount(ctx, arg.ToAccountID)
		// if err != nil {
		// 	return err
		// }

		return nil
	})

	return result, err
}
