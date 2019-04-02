package multiBalance

import (
	"log"
	"strconv"
	"sync"

	"github.com/button-tech/utils-node-tool/waves/handlers/storage"
	"github.com/imroc/req"
)

type Data struct {
	sync.Mutex
	Result map[string]string
}

func New() *Data {
	return &Data{
		Result: make(map[string]string),
	}
}

func (ds *Data) set(key string, value string) {
	ds.Result[key] = value
}

func (ds *Data) Set(key string, value string) {
	ds.Lock()
	defer ds.Unlock()
	ds.set(key, value)
}

func Worker(wg *sync.WaitGroup, addr string, r *Data) {
	defer wg.Done()

	balance, err := req.Get(storage.WavesURL + "/addresses/balance/" + addr)
	if err != nil {
		log.Println(err)
		return
	}

	var data storage.BalanceData

	err = balance.ToJSON(&data)
	if err != nil {
		log.Println(err)
	}

	balanceStr := strconv.FormatInt(data.Balance, 10)

	r.Set(addr, balanceStr)
}
