package storage

import (
	"github.com/onrik/ethrpc"
	"os"
	"sync"
)

var (
	EtcURL = os.Getenv("ETC_NODE")

	EtcClient = ethrpc.New(EtcURL)
)

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

var EtcNodeAddress NodeAddr
