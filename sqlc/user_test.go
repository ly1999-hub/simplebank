package db

import (
	"context"
	"testing"

	"github.com/ly1999-hub/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		FullName:       util.RandomString(3),
		Email:          util.CreateRandomEmail(),
	}

	user, err := testQuery.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, user.Username, arg.Username)
	require.Equal(t, user.HashedPassword, arg.HashedPassword)
	require.Equal(t, user.FullName, arg.FullName)
	require.Equal(t, user.Email, arg.Email)

	require.NotZero(t, user.ChangePasswordAt)
	require.NotZero(t, user.CreatedAt)

	argAccount := CreateAccountParams{
		Owner:   user.Username,
		Balance: 0,
		Curency: util.RandomCurency(),
	}
	account, err := testQuery.CreateAccount(context.Background(), argAccount)
	require.NoError(t, err)

	require.NotEmpty(t, account)
	require.Equal(t, argAccount.Owner, account.Owner)
	require.Equal(t, argAccount.Balance, account.Balance)
	require.Equal(t, argAccount.Curency, account.Curency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}
