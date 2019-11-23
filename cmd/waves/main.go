package main

import (
	"log"
	"os"

	"github.com/button-tech/logger"
	"github.com/button-tech/utils-node-tool/cmd/waves/handlers"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func init() {
	if err := logger.InitLogger(os.Getenv("DSN")); err != nil {
		log.Fatal(err)
	}
}

func main() {

	r := routing.New()

	g := r.Group("/waves")

	g.Get("/balance/<address>", handlers.GetBalance)

	if err := fasthttp.ListenAndServe(":8080", r.HandleRequest); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
