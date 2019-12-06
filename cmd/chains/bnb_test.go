package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/imroc/req"
)

func TestBnB(t *testing.T) {
	port := 8077
	defer startServer(t, port, r).Close()

	resp, err := req.Get("http://localhost:8080/bnb/balance/bnb1jxfh2g85q3v0tdq56fnevx6xcxtcnhtsmcu64m")
	if err != nil {
		log.Println(err)
	}

	fmt.Println(resp.String())
}
