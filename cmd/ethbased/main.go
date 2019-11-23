package main

import (
	"github.com/button-tech/utils-node-tool/cmd/ethbased/handlers"
	"github.com/button-tech/utils-node-tool/logger"
	"github.com/button-tech/utils-node-tool/nodetools/storage"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
	"os"
)

func init() {
	if err := logger.InitLogger(os.Getenv("DSN")); err != nil {
		log.Fatal(err)
	}
	startChan := make(chan struct{})
	go storage.StoreEndpointsFromDB(startChan)
	go storage.SetFastestEndpoint(startChan)
}

func main() {

	r := routing.New()

	g := r.Group("/" + os.Getenv("BLOCKCHAIN"))

	g.Get("/balance/<address>", handlers.GetBalance)

	g.Get("/transactionFee", handlers.GetTxFee)

	g.Get("/gasPrice", handlers.GetGasPrice)

	g.Get("/nonce/<address>", handlers.GetNonce)

	g.Get("/tokenBalance/<smart-contract-address>/<user-address>", handlers.GetTokenBalance)

	g.Post("/estimateGas", handlers.GetEstimateGas)

	if err := fasthttp.ListenAndServe(":8080", r.HandleRequest); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
