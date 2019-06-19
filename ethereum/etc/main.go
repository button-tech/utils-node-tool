package main

import (
	"log"
	"os"

	"github.com/button-tech/utils-node-tool/ethereum/etc/handlers"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	gin.SetMode(gin.ReleaseMode)

	etc := r.Group("/etc")

	{
		etc.GET("/balance/:address", handlers.GetBalance)

		etc.GET("/transactionFee", handlers.GetTxFee)

		etc.GET("/gasPrice", handlers.GetGasPrice)

		etc.POST("/balances", handlers.GetBalances)
	}

	if err := r.Run(":8080"); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
