package db

import (
	"context"
	"testing"

	"github.com/ly1999-hub/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	account1, _ := testQuery.GetAccount(context.Background(), int64(1))
	account2, _ := testQuery.GetAccount(context.Background(), int64(2))
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQuery.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}
