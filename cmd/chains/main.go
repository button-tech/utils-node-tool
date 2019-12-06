package main

import (
	"fmt"
	"github.com/button-tech/utils-node-tool/cmd/chains/handlers"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"io"
	"log"
	"net"
	"testing"
)

var r *routing.Router

func main() {
	r = routing.New()

	r.Get("/cosmos/balance/<address>", handlers.GetCosmosBalance)
	r.Get("/waves/balance/<address>", handlers.GetWavesBalance)
	r.Get("/xlm/balance/<address>", handlers.GetXlmBalance)
	r.Get("/xrp/balance/<address>", handlers.GetXrpBalance)
	r.Get("/tezos/balance/<address>", handlers.GetTezosBalance)
	r.Get("/bnb/balance/<address>", handlers.GetBnbBalance)
	r.Get("/algorand/balance/<address>", handlers.GetAlgorandBalance)

	if err := fasthttp.ListenAndServe(":8080", r.HandleRequest); err != nil {
		log.Fatal(err)
	}
}

func startServer(t *testing.T, port int, r *routing.Router) io.Closer {
	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		t.Fatalf("cannot start tcp server on port %d: %s", port, err)
	}
	go fasthttp.Serve(ln, r.HandleRequest)
	return ln
}
