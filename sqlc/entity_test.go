package db

import (
	"context"
	"testing"

	"github.com/ly1999-hub/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntity(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entity, err := testQuery.CreateEntry(context.Background(), arg)
	require.NoError(t, err)

	require.NotEmpty(t, entity)
	require.Equal(t, arg.AccountID, entity.AccountID)
	require.Equal(t, arg.Amount, entity.Amount)

	require.NotZero(t, entity.ID)
	require.NotZero(t, entity.CreatedAt)

	return entity
}

func TestCreateEntity(t *testing.T) {
	accout := createRandomAccount(t)
	createRandomEntity(t, accout)
}
