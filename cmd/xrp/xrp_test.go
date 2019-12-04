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

func TestXRP(t *testing.T) {
	port := 8077
	defer startServer(t, port, R).Close()

	resp, err := req.Get("http://localhost:8080/xrp/balance/rMdG3ju8pgyVh29ELPWaDuA74CpWW6Fxns")
	if err != nil {
		log.Println(err)
	}

	fmt.Println(resp.String())
}

func TestSubmitXRP(t *testing.T) {
	port := 8077
	defer startServer(t, port, R).Close()

	body := []byte(`{"devnet": false,
    "txBlob":"1200002280000000240000000361D4838D7EA4C6800000000000000000000000000055534400000000004B4E9C06F24296074F7BC48F92A97916C6DC5EA968400000000000000A732103AB40A0490F9B7ED8DF29D246BF2D6269820A0EE7742ACDD457BEA7C7D0931EDB74473045022100D184EB4AE5956FF600E7536EE459345C7BBCF097A84CC61A93B9AF7197EDB98702201CEA8009B7BEEBAA2AACC0359B41C427C1C5B550A4CA4B80CF2174AF2D6D5DCE81144B4E9C06F24296074F7BC48F92A97916C6DC5EA983143E9D4A2B8AA0780F682D136F7A56D6724EF53754"
}`)

	resp, err := req.Post("http://localhost:8080/xrp/submit-tx", req.BodyJSON(body))
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
