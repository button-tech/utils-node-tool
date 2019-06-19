package main

import (
	"github.com/button-tech/utils-node-tool/otherBlockchains/xlm/handlers"
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

	xlm := r.Group("/xlm")

	xlm.GET("/balance/:address", handlers.GetBalance)

	if err := r.Run(":8080"); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
