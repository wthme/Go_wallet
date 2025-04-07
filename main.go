package main

import (
	handler "Go_wallet/handlers"
	"Go_wallet/pglogic"
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
		log.Fatalf("Can`t get table: %v\n", err)
	}

	log.Println("Wallets table ready")


	// randuuid := "80ac1149-9777-4e18-889d-0372cf95019e"
	// dbPool.Exec(context.Background() , "UPDATE wallets SET balance = 34568 WHERE id = $1" , randuuid)


	walletdb := pglogic.NewWalletdb(dbPool)
	walletHandler := handler.NewWalletHandler(*walletdb)


	router := gin.Default()

	api := router.Group("/api/v1")
	{
		api.POST("/wallets", walletHandler.HandleWalletOperation)
		api.GET("/wallets/:walletId", walletHandler.GetWalletBalance)
	}

	err = router.Run(":" + os.Getenv("SERVER_PORT")) ; if err != nil{
		log.Fatalf("Problem with starting server %v", err)
	}

}


func getDSN() string {
	return "host=" + os.Getenv("PG_HOST") +
			" user=" + os.Getenv("PG_USER") +
			" password=" + os.Getenv("PG_PASS") +
			" dbname=" + os.Getenv("PG_DB") +
			" port=" + os.Getenv("PG_PORT")
}