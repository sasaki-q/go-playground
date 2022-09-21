package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTransaction(t *testing.T) {
	store := NewStore(testDB)

	ac := createAccountByTest(t)
	ac2 := createAccountByTest(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTransactionResult)

	for i := 0; i < n; i++ {
		go func() {
			res, err := store.TransferTransaction(context.Background(), CreateTransferParams{
				FromAccountID: ac.ID,
				ToAccountID:   ac2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- res
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		res := <-results
		require.NotEmpty(t, res)

		transfer := res.Transfer
		toEntry := res.ToEntry
		fromEntry := res.FromEntry

		require.Equal(t, transfer.Amount, toEntry.Amount)
		require.Equal(t, -transfer.Amount, fromEntry.Amount)
	}

}
