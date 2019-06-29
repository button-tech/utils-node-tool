package main

import (
	"github.com/button-tech/utils-node-tool/eth_based/handlers"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"github.com/button-tech/utils-node-tool/nodes_utils/endpoints_store"
)


func init(){
	go endpoints_store.StoreEndpoints()
}

func main() {

	r := routing.New()

	g := r.Group("/" + os.Getenv("blockchain"))

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
