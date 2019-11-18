package main

import (
	"github.com/button-tech/utils-node-tool/cmd/zilliqa/handlers"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {

	r := gin.New()

	r.Use(cors.Default())

	gin.SetMode(gin.ReleaseMode)

	g := r.Group("/zilliqa")

	g.GET("/balance/:address", handlers.GetBalance)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
