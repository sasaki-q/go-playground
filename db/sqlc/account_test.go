package db

import (
	"context"
	"database/sql"
	util "dbapp/utils"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func createAccountByTest(t *testing.T) Account {
	user := createUserByTest(t)

	param := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	res, err := testQueries.CreateAccount(ctx, param)

	require.NoError(t, err)
	require.Equal(t, param.Owner, res.Owner)
	require.Equal(t, reflect.TypeOf(res.ID).Kind(), reflect.Int64)

	return res
}

func TestCreateAcoount(t *testing.T) {
	createAccountByTest(t)
}

func TestGetAccount(t *testing.T) {
	createdAccount := createAccountByTest(t)
	selectedAccount, err := testQueries.SelectAccount(ctx, createdAccount.ID)

	require.NoError(t, err)
	require.Equal(t, createdAccount, selectedAccount)
}

func TestUpdateAccount(t *testing.T) {
	selectedAccount, _ := testQueries.SelectAccount(ctx, 10)
	param := UpdateAccountParams{
		ID:      selectedAccount.ID,
		Balance: selectedAccount.Balance + 10000,
	}

	res, err := testQueries.UpdateAccount(ctx, param)

	expected := Account{
		ID:        selectedAccount.ID,
		Owner:     selectedAccount.Owner,
		Currency:  selectedAccount.Currency,
		Balance:   selectedAccount.Balance + 10000,
		CreatedAt: selectedAccount.CreatedAt,
	}

	require.NoError(t, err)
	require.Equal(t, res, expected)
}

func TestDeleteAccount(t *testing.T) {
	createdAccount := createAccountByTest(t)

	err := testQueries.DeleteAccount(ctx, createdAccount.ID)
	require.NoError(t, err)

	_, err = testQueries.SelectAccount(ctx, createdAccount.ID)
	require.Error(t, err)
	require.Equal(t, err, sql.ErrNoRows)
}
