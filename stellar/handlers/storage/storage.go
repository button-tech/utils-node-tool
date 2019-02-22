package storage

import (
	"sync"
)

type StellarBalance struct {
	Balances []struct {
		Balance             string `json:"balance"`
		Buying_liabilities  string `json:"buying_liabilities"`
		Selling_liabilities string `json:"selling_liabilities"`
		Asset_type          string `json:"asset_type"`
	}
}

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

var StellarNodeAddress NodeAddr
