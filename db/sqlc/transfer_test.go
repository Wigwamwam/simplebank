package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wigwamwam/simplebank/util"
)

func createRandomTransfer(t *testing.T, account1 Account, account2 Account) Transfer {

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

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
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	createRandomTransfer(t, account1, account2)
}

func TestGetTranfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	transfer := createRandomTransfer(t, account1, account2)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer.ID, transfer2.ID)
	require.Equal(t, transfer.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer.Amount, transfer2.Amount)

	require.WithinDuration(t, transfer.CreatedAt, transfer2.CreatedAt, time.Second)

}

func TestListTransfers(t *testing.T) {
	var listTransfers Transfer
	var account1 Account
	var account2 Account

	for i := 0; i < 10; i++ {
		account1 = createRandomAccount(t)
		account2 = createRandomAccount(t)
		listTransfers = createRandomTransfer(t, account1, account2)
	}

	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Limit:  5,
		Offset: 0,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, listTransfers.FromAccountID, transfer.FromAccountID)
		require.Equal(t, listTransfers.ToAccountID, transfer.ToAccountID)
		require.Equal(t, listTransfers.FromAccountID, transfer.FromAccountID)
	}
}
