package db

import (
	"context"
	"database/sql"
	"fmt"
	"simpledice/util"

	"github.com/google/uuid"
)

// Store provides all functions to execute db queries and transactions
type Store interface {
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, username, asset string) (CreateUserResult, error)
	StartGame(ctx context.Context, username string) (StartGameResult, error)
	FirstDiceThrow(ctx context.Context, user GetUserWithSessionRow) error
	SecondDiceThrow(ctx context.Context, user GetUserWithSessionRow) (Attempt, error)
	StartNewDiceThrow(ctx context.Context, user GetUserWithSessionRow) (Attempt, error)
	Querier
}

// Store provides all functions to execute SQL queries and transactions
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
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromWalletID int64 `json:"from_wallet_id"`
	ToWalletID   int64 `json:"to_wallet_id"`
	Amount       int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer   Transfer `json:"transfer"`
	FromWallet Wallet   `json:"from_wallet"`
	ToWallet   Wallet   `json:"to_wallet"`
	FromEntry  Entry    `json:"from_entry"`
	ToEntry    Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one wallet to the other.
// It creates a transfer record, add wallet entries, and update wallets balance within a single database transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromWalletID: arg.FromWalletID,
			ToWalletID:   arg.ToWalletID,
			Amount:       arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			WalletID: arg.FromWalletID,
			Amount:   -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			WalletID: arg.ToWalletID,
			Amount:   arg.Amount,
		})
		if err != nil {
			return err
		}

		// TODO: update wallets balance

		if arg.FromWalletID < arg.ToWalletID {
			result.FromWallet, result.ToWallet, err = addMoney(ctx, q, arg.FromWalletID, -arg.Amount, arg.ToWalletID, arg.Amount)
		} else {
			result.ToWallet, result.FromWallet, err = addMoney(ctx, q, arg.ToWalletID, arg.Amount, arg.FromWalletID, -arg.Amount)
		}
		return nil
	})
	return result, err
}

func addMoney(ctx context.Context, q *Queries, walletID1, amount1, walletID2, amount2 int64) (wallet1, wallet2 Wallet, err error) {
	wallet1, err = q.AddWalletBalance(ctx, AddWalletBalanceParams{
		ID:     walletID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	wallet2, err = q.AddWalletBalance(ctx, AddWalletBalanceParams{
		ID:     walletID2,
		Amount: amount2,
	})
	return
}

type CreateUserResult struct {
	User   User   `json:"user"`
	Wallet Wallet `json:"wallet"`
}

func (store *SQLStore) CreateUserTx(ctx context.Context, username, asset string) (CreateUserResult, error) {
	var result CreateUserResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.User, err = q.CreateUser(ctx, username)
		if err != nil {
			return err
		}
		wallet, err := q.CreateWallet(ctx, CreateWalletParams{
			Owner:   username,
			Balance: 0,
			Asset:   asset,
		})
		if err != nil {
			return err
		}
		result.Wallet = wallet
		return nil
	})

	return result, err
}

type StartGameResult struct {
	Balance          string `json:"balance"`
	SessionID        int64  `json:"session_id"`
	CurrentAttemptID string `json:"attempt_id"`
	Asset            string `json:"asset"`
}

func (store *SQLStore) StartGame(ctx context.Context, username string) (StartGameResult, error) {
	var result StartGameResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		attemptID := uuid.New().String()
		usersession, err := q.GetUserWithSession(ctx, username)

		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if err != nil && err == sql.ErrNoRows {
			newUserSession, err := q.CreateSession(ctx, CreateSessionParams{
				IsActive: true,
				Owner:    username,
			})

			if err != nil {
				return nil
			}
			userwallet, err := q.GetWalletByUsername(ctx, username)

			if err != nil {
				return err
			}
			if userwallet.Balance-20 < 0 {
				return ErrInsufficientFund
			}
			txnRes, err := store.TransferTx(ctx, TransferTxParams{
				FromWalletID: userwallet.ID,
				ToWalletID:   util.HOLDING_ACCOUNT_ID,
				Amount:       util.COMMITMENT_FEE,
			})
			if err != nil {
				return err
			}

			userattempt, err := q.CreateAttempt(ctx, CreateAttemptParams{
				ID:           attemptID,
				SessionID:    newUserSession.ID,
				TargetNumber: int16(util.RandomInt(util.MIN_FOR_TARGET_NUM, util.MAX_FOR_TARGET_NUM)),
			})
			if err != nil {
				return err
			}
			newUserSession, err = q.UpdateSessionCurrentAttemptId(ctx, UpdateSessionCurrentAttemptIdParams{
				ID: newUserSession.ID,
				CurrentAttemptID: sql.NullString{
					String: userattempt.ID,
					Valid:  true,
				},
			})
			if err != nil {
				return err
			}
			result.Asset = txnRes.ToWallet.Asset
			result.Balance = fmt.Sprintf("%d", txnRes.FromWallet.Balance)
			result.SessionID = newUserSession.ID
			result.CurrentAttemptID = attemptID

			return nil
		}
		if usersession.IsActive {
			return ErrActiveSession
		}
		return nil
	})
	return result, err
}

func (store *SQLStore) FirstDiceThrow(ctx context.Context, user GetUserWithSessionRow) error {
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		userwallet, err := q.GetWalletByUsername(ctx, user.Owner)
		if err != nil {
			return err
		}
		if userwallet.Balance-5 < 0 {
			return ErrInsufficientFund
		}

		_, err = store.TransferTx(ctx, TransferTxParams{
			FromWalletID: userwallet.ID,
			ToWalletID:   util.HOLDING_ACCOUNT_ID,
			Amount:       util.FIRST_DICE_THROW_FEE,
		})

		if err != nil {
			return err
		}
		_, err = q.UpdateValueofFirstDiceThrown(ctx, UpdateValueofFirstDiceThrownParams{
			ID:                  user.CurrentAttemptID.String,
			FirstDiceThrowValue: int16(util.RandomInt(1, 6)),
		})

		if err != nil {
			return err
		}
		return nil
	})
	return err
}
func (store *SQLStore) SecondDiceThrow(ctx context.Context, user GetUserWithSessionRow) (Attempt, error) {
	var userattempt Attempt
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		userattempt, err = q.UpdateValueofSecondDiceThrown(ctx, UpdateValueofSecondDiceThrownParams{
			ID:                   user.CurrentAttemptID.String,
			SecondDiceThrowValue: int16(util.RandomInt(1, 6)),
		})
		if err != nil {
			return err
		}
		userwallet, err := q.GetWalletByUsername(ctx, user.Owner)
		if err != nil {
			return err
		}
		if userattempt.FirstDiceThrowValue+userattempt.SecondDiceThrowValue == userattempt.TargetNumber {
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromWalletID: util.HOLDING_ACCOUNT_ID,
				ToWalletID:   userwallet.ID,
				Amount:       util.PRIZE_AMOUNT,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return userattempt, err
}

func (store *SQLStore) StartNewDiceThrow(ctx context.Context, user GetUserWithSessionRow) (Attempt, error) {
	var userattempt Attempt
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		userwallet, err := q.GetWalletByUsername(ctx, user.Owner)
		if err != nil {
			return err
		}
		if userwallet.Balance-5 < 0 {
			return ErrInsufficientFund
		}
		attemptID := uuid.New().String()
		userattempt, err := q.CreateAttempt(ctx, CreateAttemptParams{
			ID:                  attemptID,
			SessionID:           user.ID,
			TargetNumber:        int16(util.RandomInt(util.MIN_FOR_TARGET_NUM, util.MAX_FOR_TARGET_NUM)),
			FirstDiceThrowValue: int16(util.RandomInt(1, 6)),
			NumOfDiceThrow:      1,
		})
		if err != nil {
			return err
		}
		_, err = q.UpdateSessionCurrentAttemptId(ctx, UpdateSessionCurrentAttemptIdParams{
			ID: user.ID,
			CurrentAttemptID: sql.NullString{
				String: userattempt.ID,
				Valid:  true,
			},
		})

		_, err = store.TransferTx(ctx, TransferTxParams{
			FromWalletID: userwallet.ID,
			ToWalletID:   util.HOLDING_ACCOUNT_ID,
			Amount:       util.FIRST_DICE_THROW_FEE,
		})

		if err != nil {
			return err
		}
		return nil
	})
	return userattempt, err
}
