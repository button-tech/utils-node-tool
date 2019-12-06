package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/imroc/req"
)

func TestTezos(t *testing.T) {
	port := 8077
	defer startServer(t, port, r).Close()

	resp, err := req.Get("http://localhost:8080/tezos/balance/dn1RwYfk5Mgd43xzbcD5pvDmHUsYuX4Bmbjc")
	if err != nil {
		log.Println(err)
	}

	fmt.Println(resp.String())
}
