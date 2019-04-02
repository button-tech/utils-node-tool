package multiBalance

import (
	"sync"
	"github.com/button-tech/utils-node-tool/db"
	"github.com/onrik/ethrpc"
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

	endPoint, err := db.GetEndpoint("etc")
	if err != nil{
		log.Println(err)
		return
	}

	etcClient := ethrpc.New(endPoint)


	balance, err := etcClient.EthGetBalance(addr, "latest")
	if err != nil {
		log.Println(err)
		return
	}

	r.Set(addr, balance.String())
}
