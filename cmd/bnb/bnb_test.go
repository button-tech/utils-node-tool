package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"testing"

	"github.com/imroc/req"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func TestBnB(t *testing.T) {
	port := 8077
	defer startServer(t, port, R).Close()

	resp, err := req.Get("http://localhost:8080/bnb/balance/bnb1jxfh2g85q3v0tdq56fnevx6xcxtcnhtsmcu64m")
	if err != nil {
		log.Println(err)
	}

	fmt.Println(resp.String())
}

func startServer(t *testing.T, port int, r *routing.Router) io.Closer {
	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		t.Fatalf("cannot start tcp server on port %d: %s", port, err)
	}
	go fasthttp.Serve(ln, r.HandleRequest)
	return ln
}
