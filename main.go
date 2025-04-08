package main

import (
	"Go_wallet/handlers"
	"Go_wallet/pglogic"
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading config.env file: %v", err)
	}

	dbPool, err := pgxpool.New(context.Background(), "postgresql://postgres:postgres@valet:5432/postgres")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close()

	err = pglogic.InitDB(dbPool)
	if err != nil {
		log.Fatalf("Can`t get table: %v\n", err)
	}
	log.Println("Wallets table ready")

	walletdb := pglogic.NewWalletdb(dbPool)
	walletHandler := handlers.NewWalletHandler(*walletdb)
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		api.POST("/wallets", walletHandler.HandleWalletOperation)
		api.GET("/wallets/:walletId", walletHandler.GetWalletBalance)
	}

	err = router.Run(":" + os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatalf("Problem with starting server %v", err)
	}
}
