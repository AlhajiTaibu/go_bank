package simplebank

import (
	"context"
	"database/sql"
	"fmt"
)


type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store{
	return &Store{
		db: db,
		Queries: New(db),
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error{
	tx, err := s.db.BeginTx(ctx, nil)  // Begin database transaction

	if err != nil{
		return err
	}

	q := New(tx)
	err = fn(q)

	if err !=nil{
		if rbErr := tx.Rollback(); rbErr!=nil{  // Rollback database transaction
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()  // Commit database transaction
}

type TransferMoneyTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount sql.NullInt64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

var txKey = struct{}{}

func (s *Store) TransferMoneyTx(ctx context.Context, arg TransferMoneyTxParams) (TransferTxResult, error){
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		txName := ctx.Value(txKey)
		var err error
		fmt.Println(txName, "create transfer" )
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		
		if err != nil{
			return err
		}
		fmt.Println(txName, "create from_entry" )
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: sql.NullInt64{Int64:-arg.Amount.Int64, Valid:true},
		})

		if err != nil{
			return err
		}

		fmt.Println(txName, "create to_entry" )
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: sql.NullInt64{Int64:arg.Amount.Int64, Valid:true},
		})

		if err != nil{
			return err
		}

		if arg.FromAccountID < arg.ToAccountID{
			result.FromAccount, result.ToAccount, err = transferMoney(ctx, q, arg.FromAccountID, -arg.Amount.Int64, arg.ToAccountID, arg.Amount.Int64)
		

		if err != nil{
			return err
		}

		}else{
			result.FromAccount, result.ToAccount, err = transferMoney(ctx, q, arg.ToAccountID, arg.Amount.Int64, arg.FromAccountID, -arg.Amount.Int64)


		if err != nil{
			return err
		}
		}
		return nil
	})

	return result, err
}

func transferMoney(
	ctx context.Context,
	q *Queries,
	FromAccountID int64,
	fromAccountAmount int64,
	ToAccountID int64,
	toAccountAmount int64,
)(Account, Account, error){
	fromAccount, err := q.GetAccountForUpdate(ctx, FromAccountID)

	if err != nil{
		return Account{}, Account{}, err
	}
		
	FromAccount, err := q.UpdateAccount(ctx, UpdateAccountParams{
			ID: FromAccountID,
			Balance: fromAccount.Balance + fromAccountAmount,
		})

	if err != nil{
		return FromAccount, Account{}, err
	}

	
	toAccount, err := q.GetAccountForUpdate(ctx, ToAccountID)

	if err != nil{
		return FromAccount, Account{}, err
	}
	
	ToAccount, err := q.UpdateAccount(ctx, UpdateAccountParams{
		ID: ToAccountID,
		Balance: toAccount.Balance + toAccountAmount,
	})

	if err != nil{
		return FromAccount, ToAccount, err
	}
	return FromAccount, ToAccount, nil
}