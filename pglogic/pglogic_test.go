package pglogic

import (
	"Go_wallet/model"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestGetWallet(t *testing.T) {

	dbPool, _  := pgxpool.New(context.Background(), "postgresql://postgres:12345@valet:5432/postgres")

	wdb := NewWalletdb(dbPool)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	walletID := uuid.New()

	t.Run("Test Init Wallet", func(t *testing.T) {

		_, err  := wdb.GetWallet(ctx , walletID)

		assert.NoError(t, err)
	})

	t.Run("Test Deposit", func(t *testing.T) {
		req := model.WalletOperation{
			WalletID:      walletID,
			OperationType: model.DEPOSIT,
			Amount:        1488,
		}
		err := wdb.ProcessOperation(ctx , req)
		assert.NoError(t, err)
		res ,_:= wdb.GetWalletBalance(ctx , req.WalletID)
		assert.Equal(t, float64(1488), res)
	})



	t.Run("Test Concurrent Transactions", func(t *testing.T) {
		const concurrentRequests = 100
		results := make(chan error, concurrentRequests)
		walletID = uuid.New()
		wdb.CreateWallet(ctx,walletID)


		for i := 0; i < concurrentRequests; i++ {
			go func() {
				req := model.WalletOperation{
					WalletID:      walletID,
					OperationType: model.DEPOSIT,
					Amount:        1,
				}

				err := wdb.ProcessOperation(ctx, req)
				results <- err
			}()
		}

		for i := 0; i < concurrentRequests; i++ {
			err := <-results
			assert.NoError(t, err)
		}

		wallet, err := wdb.GetWallet(ctx, walletID)
		assert.NoError(t, err)
		assert.Equal(t, float64(concurrentRequests), wallet.Balance)
	})

	t.Run("Test Withdraw", func(t *testing.T) {

		req := model.WalletOperation{
			WalletID:      walletID,
			OperationType: model.WITHDRAW,
			Amount:        20,
		}
		err := wdb.ProcessOperation(ctx , req)
		assert.NoError(t, err)
		res ,_:= wdb.GetWalletBalance(ctx , req.WalletID)
		assert.Equal(t, float64(80), res)
	})

}