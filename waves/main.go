package main

import (
	_ "github.com/button-tech/utils-node-tool/waves/docs"
	"github.com/button-tech/utils-node-tool/waves/handlers"
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

	r.GET("/waves/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/waves/balance/:address", handlers.GetBalance)

	r.Run(":8080")
}
