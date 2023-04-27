package db

import (
	"context"
	"simpledice/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, account1, account2 Wallet) Transfer {
	arg := CreateTransferParams{
		FromWalletID: account1.ID,
		ToWalletID:   account2.ID,
		Amount:       util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromWalletID, transfer.FromWalletID)
	require.Equal(t, arg.ToWalletID, transfer.ToWalletID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	account1 := createRandomWallet(t)
	account2 := createRandomWallet(t)
	createRandomTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomWallet(t)
	account2 := createRandomWallet(t)
	transfer1 := createRandomTransfer(t, account1, account2)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromWalletID, transfer2.FromWalletID)
	require.Equal(t, transfer1.ToWalletID, transfer2.ToWalletID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestListTransfer(t *testing.T) {
	account1 := createRandomWallet(t)
	account2 := createRandomWallet(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, account1, account2)
		createRandomTransfer(t, account2, account1)
	}

	arg := ListTransfersParams{
		FromWalletID: account1.ID,
		ToWalletID:   account1.ID,
		Limit:        5,
		Offset:       5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromWalletID == account1.ID || transfer.ToWalletID == account1.ID)
	}
}
