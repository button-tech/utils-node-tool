package main

import (
	"github.com/button-tech/utils-node-tool/eth-based/handlers"
	"github.com/button-tech/utils-node-tool/utils-for-endpoints/estorage"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
	"os"
)

func init() {
	go estorage.StoreEndpointsFromDB()
	go estorage.SetFastestEndpoint()
}

func main() {

	r := routing.New()

	g := r.Group("/" + os.Getenv("BLOCKCHAIN"))

	g.Get("/balance/<address>", handlers.GetBalance)

	g.Get("/transactionFee", handlers.GetTxFee)

	g.Get("/gasPrice", handlers.GetGasPrice)

	g.Get("/tokenBalance/<smart-contract-address>/<user-address>", handlers.GetTokenBalance)

	g.Post("/estimateGas", handlers.GetEstimateGas)

	if err := fasthttp.ListenAndServe(":8080", r.HandleRequest); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
