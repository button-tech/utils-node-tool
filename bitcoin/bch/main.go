package main

import (
	"github.com/button-tech/utils-node-tool/bitcoin/bch/handlers"
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

	bch := r.Group("/bch")

	{
		bch.GET("/balance/:address", handlers.GetBalance)

		bch.GET("/utxo/:address", handlers.GetUTXO)

		bch.GET("/transactionFee", handlers.GetTxFee)

		bch.POST("/balances", handlers.GetBalances)
	}

	if err := r.Run(":8080"); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
