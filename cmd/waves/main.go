package main

import (
	"log"
	"os"

	"github.com/button-tech/utils-node-tool/cmd/waves/handlers"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func main() {

	r := routing.New()

	g := r.Group("/waves")

	g.Get("/balance/<address>", handlers.GetBalance)

	// g.Post("/balances", handlers.GetBalances)

	if err := fasthttp.ListenAndServe(":8080", r.HandleRequest); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
