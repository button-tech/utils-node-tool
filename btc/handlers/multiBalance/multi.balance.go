package multiBalance

import (
	"github.com/button-tech/utils-node-tool/shared/db"
	"github.com/imroc/req"
	"log"
	"sync"
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

	endPoint, err := db.GetEndpoint("btc")
	if err != nil {
		log.Println(err)
		return
	}

	balance, err := req.Get(endPoint + "/addr/" + addr + "/balance")
	if err != nil {
		log.Println(err)
		return
	}

	balanceStr, err := balance.ToString()
	if err != nil {
		log.Println(err)
		return
	}

	r.Set(addr, balanceStr)
}
