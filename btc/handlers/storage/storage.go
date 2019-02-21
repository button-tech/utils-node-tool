package storage

import (
	"os"
	"sync"
)

var BtcURL = os.Getenv("BTC_NODE")

type NodeAddr struct {
	sync.Mutex
	Address string
}

func (na *NodeAddr) set(value string) {
	na.Address = value
}

func (ds *NodeAddr) Set(value string) {
	ds.Lock()
	defer ds.Unlock()
	ds.set(value)
}

var BtcNodeAddress NodeAddr
