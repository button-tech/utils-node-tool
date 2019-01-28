package multiBalance

import (
	"fmt"
	"github.com/button-tech/utils-node-tool/etc/storage"
	"math"
	"strconv"
	"sync"
)

type Balances struct {
	sync.Mutex
	Result []map[string]float64 // eth -> ["account":balance]
}

func (ds *Balances) set(key string, value float64) {
	ds.Result = append(ds.Result, map[string]float64{key: value})
}

func (ds *Balances) Set(key string, value float64) {
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
	floatBalance, _ := strconv.ParseFloat(balance.String(), 64)
	ethBalance := floatBalance / math.Pow(10, 18)
	r.Set(addr, ethBalance)
}
