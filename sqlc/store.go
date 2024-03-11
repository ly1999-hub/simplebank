package db

import (
	"context"
	"database/sql"
	"fmt"
)

// TransferTxParams containt input parameters of transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is result transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"Transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntity  Entry    `json:"from_entity"`
	ToEntity    Entry    `json:"to_entity"`
}

var txKey = struct{}{}

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		return err

	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v ,RollBack Error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// TransferTx perform a money transfer from one account to the other
// It creates transfer record, add account entries, and update account balance within a single database transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var (
		result TransferTxResult
	)
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		txName := ctx.Value(txKey)
		fmt.Println(txName)
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntity, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntity, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		// Todo update balance
		account1, err := q.GetAccountForUpdate(context.Background(), arg.FromAccountID)
		if err != nil {
			return err
		}

		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.FromAccountID,
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}

		// Todo update balance
		account2, err := q.GetAccountForUpdate(context.Background(), arg.ToAccountID)
		if err != nil {
			return err
		}

		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.ToAccountID,
			Balance: account2.Balance + arg.Amount,
		})
		if err != nil {
			return err
		}

		if err != nil {
			return err
		}
		return nil
	})
	return result, err
}
