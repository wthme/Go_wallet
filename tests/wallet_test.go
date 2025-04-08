package tests

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"Go_wallet/model"
	"Go_wallet/pglogic"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)




func getDSN() string {
	return "host=" + os.Getenv("PG_HOST") +
			" user=" + os.Getenv("PG_USER") +
			" password=" + os.Getenv("PG_PASS") +
			" dbname=" + os.Getenv("PG_DB") +
			" port=" + os.Getenv("PG_PORT")
}

func TestWalletRepository(t *testing.T) {

	if err := godotenv.Load("../config.env"); err != nil {
		log.Fatalf("Error loading config.env file: %v", err)
	}


	dbPool, err := pgxpool.New(context.Background(), getDSN())
	if err != nil {
		log.Fatalf("can`t  make test connection: %v\n", err)
	}
	defer dbPool.Close()



	repo := pglogic.NewWalletdb(dbPool)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	walletID := uuid.New()

	t.Run("Test Init Wallet", func(t *testing.T) {
		err := repo.CreateWallet(ctx, walletID)
		assert.NoError(t, err)
	})

	t.Run("Test Get Wallet", func(t *testing.T) {
		wallet, err := repo.GetWallet(ctx, walletID)
		assert.NoError(t, err)
		assert.Equal(t, walletID, wallet.ID)
		assert.Equal(t, float64(0), wallet.Balance)
	})

	t.Run("Test Deposit", func(t *testing.T) {
		req := model.WalletOperation{
			WalletID:      walletID,
			OperationType: model.DEPOSIT,
			Amount:        1000,
		}

		err := repo.ProcessOperation(ctx, req)
		assert.NoError(t, err)
		// assert.Equal(t, int64(1000), wallet.Balance)

		// Verify balance
		wallet, err := repo.GetWallet(ctx, walletID)
		assert.NoError(t, err)
		assert.Equal(t, float64(1000), wallet.Balance)
	})

	t.Run("Test Withdraw", func(t *testing.T) {
		req := model.WalletOperation{
			WalletID:      walletID,
			OperationType: model.WITHDRAW,
			Amount:        500,
		}

		err := repo.ProcessOperation(ctx, req)
		assert.NoError(t, err)
		// assert.Equal(t, int64(500), req.WalletID.)

		// Verify balance
		// wallet, err = repo.GetWallet(ctx, walletID)
		// assert.NoError(t, err)
		// assert.Equal(t, int64(500), wallet.Balance)
	})

	// t.Run("Test Insufficient Balance", func(t *testing.T) {
	// 	req := model.WalletOperation{
	// 		WalletID:      walletID,
	// 		OperationType: model.WITHDRAW,
	// 		Amount:        1000,
	// 	}

	// 	// err := repo.ProcessTransaction(ctx, req)
	// 	// assert.ErrorIs(t, err, repository.ErrInsufficientBalance)

	// 	// Verify balance didn't change
	// 	wallet, err := repo.GetWallet(ctx, walletID)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, int64(500), wallet.Balance)
	// })

	t.Run("Test Concurrent Transactions", func(t *testing.T) {

		// walletID := uuid.New()



		const concurrentRequests = 100


		results := make(chan error, concurrentRequests)



		for i := 0; i < concurrentRequests; i++ {
			go func() {
				req := model.WalletOperation{
					WalletID:      walletID,
					OperationType: model.DEPOSIT,
					Amount:        1,
				}

				err := repo.ProcessOperation(ctx, req)
				results <- err
			}()
		}

		for i := 0; i < concurrentRequests; i++ {
			err := <-results
			assert.NoError(t, err)
		}

		// Verify final balance
		wallet, err := repo.GetWallet(ctx, walletID)
		assert.NoError(t, err)
		assert.Equal(t, float64(500+concurrentRequests), wallet.Balance)
	})
}







// func SetupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
// 	t.Helper()

// 	pool, err := dockertest.NewPool("")
// 	if err != nil {
// 		t.Fatalf("Could not connect to docker: %s", err)
// 	}

// 	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
// 		Repository: "postgres",
// 		Tag:        "15-alpine",
// 		Env: []string{
// 			"POSTGRES_USER=test_user",
// 			"POSTGRES_PASSWORD=test_password",
// 			"POSTGRES_DB=test_db",
// 			"listen_addresses = '*'",
// 		},
// 	}, func(config *docker.HostConfig) {
// 		config.AutoRemove = true
// 		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
// 	})
// 	if err != nil {
// 		t.Fatalf("Could not start resource: %s", err)
// 	}

// 	hostAndPort := resource.GetHostPort("5432/tcp")
// 	databaseUrl := fmt.Sprintf("postgres://test_user:test_password@%s/test_db?sslmode=disable", hostAndPort)

// 	// Exponential backoff-retry
// 	var dbPool *pgxpool.Pool
// 	if err := pool.Retry(func() error {
// 		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 		defer cancel()

// 		var err error
// 		dbPool, err = pgxpool.New(ctx, databaseUrl)
// 		if err != nil {
// 			return err
// 		}
// 		return dbPool.Ping(ctx)
// 	}); err != nil {
// 		t.Fatalf("Could not connect to docker: %s", err)
// 	}

// 	cleanup := func() {
// 		dbPool.Close()
// 		pool.Purge(resource)
// 	}

// 	return dbPool, cleanup
// }