package storage

import (
	"github.com/onrik/ethrpc"
	"os"
	"sync"
)

var (
	EthURL = os.Getenv("ETH_NODE")

	EthClient = ethrpc.New(EthURL)

	TokensAddresses = map[string]string{

		"bix":  "0xb3104b4b9da82025e8b9f8fb28b3553ce2f67069",
		"btm":  "0xcb97e65f07da24d46bcdd078ebebd7c6e6e3d750",
		"omg":  "0xd26114cd6EE289AccF82350c8d8487fedB8A0C07",
		"elf":  "0xbf2179859fc6d5bee9bf9158632dc51678a4100e",
		"bnb":  "0xB8c77482e45F1F44dE1745F52C74426C631bDD52",
		"tusd": "0x8dd5fbce2f6a956c3022ba3663759011dd51e73e",
		"knc":  "0xdd974d5c2e2928dea5f71b9825b8b646686bd200",
		"zrx":  "0xe41d2489571d322189246dafa5ebde1f4699f498",
		"rep":  "0x1985365e9f78359a9B6AD760e32412f4a445E862",
		"gnt":  "0xa74476443119A942dE498590Fe1f2454d7D4aC0d",
	}

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

var EthNodeAddress NodeAddr
