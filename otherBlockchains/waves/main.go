package main

import (
	"github.com/button-tech/utils-node-tool/otherBlockchains/waves/handlers"
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

	waves := r.Group("/waves")

	{
		waves.GET("/balance/:address", handlers.GetBalance)

		waves.GET("/transactionFee", handlers.GetTxFee)

		waves.POST("/balances", handlers.GetBalances)
	}

	if err := r.Run(":8080"); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
