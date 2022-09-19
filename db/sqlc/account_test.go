package db

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    "tom",
		Balance:  100,
		Currency: "USD",
	}

	res, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Equal(t, arg.Owner, res.Owner)
	require.Equal(t, reflect.TypeOf(res.ID).Kind(), reflect.Int64)
}
