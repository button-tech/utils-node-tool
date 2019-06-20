package main

import (
	"github.com/button-tech/utils-node-tool/utxoBased/handlers"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
	"os"
)

func main() {

	r := routing.New()

	g := r.Group("/" + os.Getenv("blockchain"))

	g.Get("/balance/<address>", handlers.GetBalance)

	g.Get("/utxo/<address>", handlers.GetUTXO)

	g.Post("/balances", handlers.GetBalances)

	if err := fasthttp.ListenAndServe(":8080", r.HandleRequest); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
