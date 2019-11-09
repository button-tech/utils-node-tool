package main

import (
	"github.com/button-tech/utils-node-tool/cmd/tezos/handlers"
	"log"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func main() {
	r := routing.New()
	g := r.Group("/tezos")
	g.Get("/balance/<address>", handlers.GetBalance)

	if err := fasthttp.ListenAndServe(":8080", r.HandleRequest); err != nil {
		log.Fatal(err)
	}
}
