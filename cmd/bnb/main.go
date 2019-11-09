package main

import (
	"log"

	"github.com/button-tech/utils-node-tool/cmd/bnb/handlers"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

var R *routing.Router

func main() {
	R := routing.New()
	g := R.Group("/bnb")
	g.Get("/balance/<address>", handlers.GetBalance)

	if err := fasthttp.ListenAndServe(":8080", R.HandleRequest); err != nil {
		log.Fatal(err)
	}
}
