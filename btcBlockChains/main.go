package main

import (
	"github.com/button-tech/utils-node-tool/btcBlockChains/handlers"
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

	group := r.Group("/" + os.Getenv("blockChain"))

	{
		group.GET("/balance/:address", handlers.GetBalance)

		group.GET("/utxo/:address", handlers.GetUTXO)

		//btc.POST("/balances", handlers.GetBalances)
	}

	if err := r.Run(":8080"); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
