package main

import (
	"github.com/button-tech/utils-node-tool/waves/handlers"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
	"os"
)

func main() {

	r := routing.New()

	waves := r.Group("/waves")

	waves.Get("/balance/<address>", handlers.GetBalance)

	waves.Post("/balances", handlers.GetBalances)

	if err := fasthttp.ListenAndServe(":8080", r.HandleRequest); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
