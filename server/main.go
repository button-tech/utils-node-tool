package main

import (
	_ "./docs"
	"github.com/button-tech/utils-node-tool/server/handlers"
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

	// @host localhost:8080
	// @BasePath /

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.Default())
	gin.SetMode(gin.ReleaseMode)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/balance/:nodeType/:address", handlers.GetBalance)

	r.GET("/transactionFee/:nodeType", handlers.GetTxFee)

	r.GET("/gasPrice/:nodeType", handlers.GetGasPrice)

	r.GET("/tokenBalance/:token/:address", handlers.GetTokenBalance)

	r.POST("/sendTx/:nodeType", handlers.SendTX)

	r.Run(":8080")
}
