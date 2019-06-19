package main

import (
	"github.com/button-tech/utils-node-tool/ethBlockChains/handlers"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
	"os"
)

func main() {

	r := routing.New()

	eth := r.Group("/" + os.Getenv("blockChain"))

	eth.Get("/balance/<address>", handlers.GetBalance)

	eth.Get("/transactionFee", handlers.GetTxFee)

	eth.Get("/gasPrice", handlers.GetGasPrice)

	eth.Get("/tokenBalance/<sc-address>/<address>", handlers.GetTokenBalance)

	eth.Post("/estimateGas", handlers.GetEstimateGas)

	//eth.POST("/balances", handlers.GetBalances)
	//
	//eth.POST("/tokenBalances", handlers.GetTokenBalances)

	if err := fasthttp.ListenAndServe(":8080", r.HandleRequest); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
