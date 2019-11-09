package main

import (
	"log"

	"github.com/button-tech/utils-node-tool/cmd/xrp/handlers"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func main() {
	r := routing.New()
	g := r.Group("/xrp")
	g.Get("/balance/<address>", handlers.GetBalance)

	if err := fasthttp.ListenAndServe(":8080", r.HandleRequest); err != nil {
		log.Fatal(err)
	}
}
