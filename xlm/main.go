package main

import (
	"github.com/button-tech/utils-node-tool/xlm/handlers"
	"log"
	"os"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func main() {

	r := routing.New()

	xlm := r.Group("/xlm")

	xlm.Get("/balance/<address>", handlers.GetBalance)

	if err := fasthttp.ListenAndServe(":8080", r.HandleRequest); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
