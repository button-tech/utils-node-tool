package main

import (
	"github.com/button-tech/utils-node-tool/btcBlockChains/handlers"
	"log"
	"os"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func main() {

	r := routing.New()

	btc := r.Group("/" + os.Getenv("blockChain"))

	btc.Get("/balance/<address>", handlers.GetBalance)

	btc.Get("/utxo/<address>", handlers.GetUTXO)

	if err := fasthttp.ListenAndServe(":8080", r.HandleRequest); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}