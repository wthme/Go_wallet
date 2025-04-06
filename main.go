package main

import (
	handler "Go_wallet/handlers"
	"Go_wallet/pglogic"
	"Go_wallet/service"
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load("config.env"); err != nil {
		log.Fatalf("Error loading config.env file: %v", err)
	}


	dbPool, err := pgxpool.New(context.Background(), getDSN())

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close()


	err = pglogic.InitDB(dbPool)

	if err!= nil{
		log.Fatal("Can`t get table: %v\n", err)
	}

	log.Println("Wallets table ready")

	// Initialize repository, service and handler
	wallet := pglogic.NewWallet(dbPool)
	walletService := service.NewWalletService(wallet)
	walletHandler := handler.NewWalletHandler(walletService)

	// Set up Gin router
	router := gin.Default()

	// API routes
	api := router.Group("/api/v1")
	{
		api.POST("/wallet", walletHandler.HandleWalletOperation)
		api.GET("/wallets/:walletId", walletHandler.GetWalletBalance)
	}

	// Start server
	port := os.Getenv("SERVER_PORT")
	log.Printf("Server running on port %s", port)

	log.Fatal(router.Run(":" + port))
}



func getDSN() string {
	return "host=" + os.Getenv("PG_HOST") +
			" user=" + os.Getenv("PG_USER") +
			" password=" + os.Getenv("PG_PASS") +
			" dbname=" + os.Getenv("PG_DB") +
			" port=" + os.Getenv("PG_PORT") +
			" sslmode=disable TimeZone=UTC"
}