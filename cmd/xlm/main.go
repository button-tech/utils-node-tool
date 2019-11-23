package main

import (
	"github.com/button-tech/logger"
	"github.com/button-tech/utils-node-tool/cmd/xlm/handlers"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
	"os"
)

func init() {
	if err := logger.InitLogger(os.Getenv("DSN")); err != nil {
		log.Fatal(err)
	}
}

func main() {

	r := routing.New()

	g := r.Group("/xlm")

	g.Get("/balance/<address>", handlers.GetBalance)

	if err := fasthttp.ListenAndServe(":8080", r.HandleRequest); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
