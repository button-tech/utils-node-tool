package main

import (
	"github.com/button-tech/logger"
	"github.com/button-tech/utils-node-tool/cmd/zilliqa/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func init() {
	if err := logger.InitLogger(os.Getenv("DSN")); err != nil {
		log.Fatal(err)
	}
}

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
