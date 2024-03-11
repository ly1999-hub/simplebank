package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)
	n := 2
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		txName, _ := fmt.Printf("tx:%d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			fmt.Println(result.ToEntity)
			fmt.Println(result.FromAccount)
			fmt.Println(result.ToAccount)
			fmt.Println(result.FromEntity)
			fmt.Println(result.Transfer)
			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		fmt.Println("check-transfer")

		// check transfer
		transfer := result.Transfer
		fmt.Println(transfer, account1, account2)
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fmt.Println("check-entity")
		// check entries
		fromEntity := result.FromEntity
		require.NotZero(t, fromEntity)
		require.Equal(t, account1.ID, fromEntity.AccountID)
		require.Equal(t, -amount, fromEntity.Amount)
		require.NotZero(t, fromEntity.ID)
		require.NotZero(t, fromEntity.CreatedAt)
		_, err = store.GetEntry(context.Background(), fromEntity.ID)
		require.NoError(t, err)

		toEntity := result.ToEntity
		require.NotZero(t, toEntity)
		require.Equal(t, account2.ID, toEntity.AccountID)
		require.Equal(t, amount, toEntity.Amount)
		require.NotZero(t, toEntity.ID)
		require.NotZero(t, toEntity.CreatedAt)
		_, err = store.GetEntry(context.Background(), toEntity.ID)
		require.NoError(t, err)

		fmt.Println("check-acount")
		// check account
		fromAcount := result.FromAccount
		fmt.Println("fromAcount", fromAcount)
		require.NotEmpty(t, fromAcount)
		require.Equal(t, account1.ID, fromAcount.ID)

		toAcount := result.ToAccount
		fmt.Println("toAcount.ID", toAcount.ID, account2.ID)

		require.Equal(t, account2.ID, toAcount.ID)

		fmt.Println("tx: ", fromAcount.Balance, ",", toAcount.Balance)

		// check balance
		diff1 := account1.Balance - fromAcount.Balance
		diff2 := toAcount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccount1, err := testQuery.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQuery.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println(">> After update", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)

}
