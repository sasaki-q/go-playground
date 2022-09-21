package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type Store struct {
	*Queries
	db *sql.DB
}

// NewStore create Store instance
func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (store *Store) execTransaction(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		log.Fatal("ERROR: transaction error === ", err)
		return err
	}

	fmt.Print("DEBUG success create transction instance === ", tx)

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("ERROR: transaction error: %v, rollback error: %v", err, rbErr)
		}

		return fmt.Errorf("ERROR: transaction error: %v", err)
	}

	return tx.Commit()
}

type TransferTransactionResult struct {
	Transfer      Transfer `json:"transfer"`
	FromAccountID int64    `json:"from_account_id"`
	ToAccountID   int64    `json:"to_account_id"`
	FromEntry     Entry    `json:"from_entry"`
	ToEntry       Entry    `json:"to_entry"`
}

/*
	step
	1 record transfer
	2 record entry from
	3 record entry to
	4 update account balance
*/
func (store *Store) TransferTransaction(ctx context.Context, param CreateTransferParams) (TransferTransactionResult, error) {
	var res TransferTransactionResult

	err := store.execTransaction(ctx, func(q *Queries) error {
		var txerr error

		// transfer: interactive record
		res.Transfer, txerr = q.CreateTransfer(ctx, param)

		if txerr != nil {
			return txerr
		}

		//　entry record one side
		res.FromEntry, txerr = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: param.FromAccountID,
			Amount:    -param.Amount,
		})

		if txerr != nil {
			return txerr
		}

		res.ToEntry, txerr = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: param.ToAccountID,
			Amount:    param.Amount,
		})

		if txerr != nil {
			return txerr
		}

		_, txerr = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      param.ToAccountID,
			Balance: param.Amount,
		})

		if txerr != nil {
			return txerr
		}

		_, txerr = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      param.FromAccountID,
			Balance: -param.Amount,
		})

		if txerr != nil {
			return txerr
		}

		return nil
	})

	return res, err
}
