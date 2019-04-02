package main

import (
	"log"
	"os"

	_ "github.com/button-tech/utils-node-tool/ltc/docs"
	"github.com/button-tech/utils-node-tool/ltc/handlers"
	"github.com/button-tech/utils-node-tool/ltc/handlers/storage"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prazd/round-robin"
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

	// must add addresses to slice
	var LtcNodes = []string{}

	rr, err := roundrobin.New(LtcNodes)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	// Round Robin middleware
	r.Use(func(c *gin.Context) {
		storage.LtcNodeAddress.Set(rr.Next())
	})

	gin.SetMode(gin.ReleaseMode)

	r.GET("/ltc/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ltc/balance/:address", handlers.GetBalance)

	r.GET("/ltc/utxo/:address", handlers.GetUTXO)

	r.GET("/ltc/transactionFee", handlers.GetTxFee)

	r.POST("/ltc/balances", handlers.GetBalances)

	if err = r.Run(":8080"); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
