package main

import (
	"log"
	"os"

	_ "github.com/button-tech/utils-node-tool/bitcoin/btc/docs"
	"github.com/button-tech/utils-node-tool/bitcoin/btc/handlers"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func main() {

	// @title Swagger BUTTON Node API
	// @version 1.0
	// @description This is BUTTON Node API responseModels documentation

	// @contact.name API Support
	// @contact.email nk
	// ap@buttonwallet.com

	// @license.name MIT
	// @license.url https://opensource.org/licenses/MIT

	// @BasePath /

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	gin.SetMode(gin.ReleaseMode)

	btc := r.Group("/btc")

	{
		btc.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
