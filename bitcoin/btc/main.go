package main

import (
	"github.com/button-tech/utils-node-tool/bitcoin/btc/handlers"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	gin.SetMode(gin.ReleaseMode)

	btc := r.Group("/btc")

	{
		btc.GET("/balance/:address", handlers.GetBalance)

		btc.GET("/utxo/:address", handlers.GetUTXO)

		btc.GET("/transactionFee", handlers.GetTxFee)

		btc.GET("/bestTransactionFee", handlers.GetBextTxFee)

		btc.POST("/balances", handlers.GetBalances)
	}

	if err := r.Run(":8080"); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
