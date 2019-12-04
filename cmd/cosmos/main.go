package main

import (
	"github.com/button-tech/utils-node-tool/cmd/cosmos/handlers"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
)

var r *routing.Router

func main() {
	r = routing.New()
	g := r.Group("/cosmos")
	g.Get("/balance/<address>", handlers.GetBalance)

	if err := fasthttp.ListenAndServe(":8080", r.HandleRequest); err != nil {
		log.Fatal(err)
	}
}
