package main

import (
	"github.com/button-tech/utils-node-tool/cmd/utxobased/handlers"
	"github.com/button-tech/utils-node-tool/logger"
	"github.com/button-tech/utils-node-tool/nodetools/storage"
	"github.com/imroc/req"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
	"os"
)

func init() {
	if err := logger.InitLogger(os.Getenv("DSN")); err != nil {
		log.Fatal(err)
	}
	startChan := make(chan struct{})
	go storage.StoreEndpointsFromDB(startChan)
	go storage.SetFastestEndpoint(startChan)

	req.SetClient(&http.Client{})
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
