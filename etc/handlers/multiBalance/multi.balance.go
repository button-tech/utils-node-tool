package multiBalance

import (
	"fmt"
	"github.com/button-tech/utils-node-tool/etc/handlers/storage"
	"sync"
)

type Balances struct {
	sync.Mutex
	Result []map[string]string // eth -> ["account":balance]
}

func (ds *Balances) set(key string, value string) {
	ds.Result = append(ds.Result, map[string]string{key: value})
}

func (ds *Balances) Set(key string, value string) {
	ds.Lock()
	defer ds.Unlock()
	ds.set(key, value)
}

func Worker(wg *sync.WaitGroup, addr string, r *Balances) {
	defer wg.Done()

	balance, err := storage.EtcClient.EthGetBalance(addr, "latest")
	if err != nil {
		fmt.Println(err)
		return
	}

	r.Set(addr, balance.String())
}
