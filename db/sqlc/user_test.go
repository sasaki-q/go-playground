package db

import (
	util "dbapp/utils"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func createUserByTest(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	param := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	res, err := testQueries.CreateUser(ctx, param)

	require.NoError(t, err)
	require.Equal(t, param.Username, res.Username)
	require.Equal(t, reflect.TypeOf(res.Username).Kind(), reflect.String)

	return res
}

func TestCreateUser(t *testing.T) {
	createAccountByTest(t)
}

func TestGetUser(t *testing.T) {
	createdUser := createUserByTest(t)
	selectedUser, err := testQueries.SelectUser(ctx, createdUser.Username)

	require.NoError(t, err)
	require.Equal(t, createdUser, selectedUser)
}
