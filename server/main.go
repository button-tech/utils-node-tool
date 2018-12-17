package main

import (
	"./handlers"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "../docs"
)

func main() {

	// @title Swagger BUTTON Node API
	// @version 1.0
	// @description This is BUTTON Node API responseModels documentation

	// @contact.name API Support
	// @contact.email nk@buttonwallet.com

	// @license.name MIT
	// @license.url https://opensource.org/licenses/MIT

	// @host localhost:8080
	// @BasePath /v2

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.Default())
	gin.SetMode(gin.ReleaseMode)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/balance/:eth/:address", handlers.GetBalance)

	// for ETH/ETC fee
	r.GET("/fee/:eth", handlers.GetTxFee)

	// ETH/ETC gas price
	r.GET("/gasprice/:eth", handlers.GetGasPrice)

	// for ETH only
	// get ETH acc token balance
	// example - /tbalance/knc/*
	r.GET("/tbalance/:token/:address", handlers.GetTokenBalance)

	// send rawTx ETH/ETC
	// example {"raw_tx":"f86d8202b284773594008252..."}
	r.POST("/send_tx/:eth", handlers.SendTX)

	r.Run(":8080")
}
