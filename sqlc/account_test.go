package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/ly1999-hub/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:   user.Username,
		Balance: util.RandomMoney(),
		Curency: util.RandomCurency(),
	}

	account, err := testQuery.CreateAccount(context.Background(), arg)
	require.NoError(t, err)

	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Curency, account.Curency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccout(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account, err := testQuery.GetAccount(context.Background(), 1)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(account)
}

func TestListAccount(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}
	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}
	accounts, err := testQuery.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}
