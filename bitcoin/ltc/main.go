package main

import (
	"github.com/button-tech/utils-node-tool/bitcoin/ltc/handlers"
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

	ltc := r.Group("/ltc")

	{

		ltc.GET("/balance/:address", handlers.GetBalance)

		ltc.GET("/utxo/:address", handlers.GetUTXO)

		ltc.GET("/transactionFee", handlers.GetTxFee)

		ltc.POST("/balances", handlers.GetBalances)
	}

	if err := r.Run(":8080"); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
