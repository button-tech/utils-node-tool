package main

import (
	_ "./docs"
	"github.com/button-tech/utils-node-tool/etc/handlers"
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

	r.GET("/etc/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/etc/balance/:address", handlers.GetBalance)

	r.GET("/etc/transactionFee", handlers.GetTxFee)

	r.GET("/etc/gasPrice/", handlers.GetGasPrice)

	// r.POST("/eth/sendTx/", handlers.SendTX)

	r.Run(":8080")
}
