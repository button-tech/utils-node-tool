package main

import (
	"github.com/button-tech/utils-node-tool/utils-for-endpoints/storage"
	"github.com/button-tech/utils-node-tool/utxo-based/handlers"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"time"
)

func init() {
	go storage.StoreEndpointsFromDB()
	time.Sleep(time.Second * 1)
	go storage.SetFastestEndpoint()
}

func main() {

	r := routing.New()

	g := r.Group("/" + os.Getenv("BLOCKCHAIN"))

	g.Get("/balance/<address>", handlers.GetBalance)

	g.Get("/utxo/<address>", handlers.GetUtxo)

	g.Post("/balances", handlers.GetBalances)

	if err := fasthttp.ListenAndServe(":8080", r.HandleRequest); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
