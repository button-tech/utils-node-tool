package multiBalance

import (
	"fmt"
	"github.com/button-tech/utils-node-tool/eth/abi"
	"github.com/button-tech/utils-node-tool/eth/handlers/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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
	balance, err := storage.EthClient.EthGetBalance(addr, "latest")
	if err != nil {
		fmt.Println(err)
		return
	}

	r.Set(addr, balance.String())
}

func TokenWorker(wg *sync.WaitGroup, address string, smartContractAddress string, r *Data) {

	defer wg.Done()

	ethClient, err := ethclient.Dial(storage.EthURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	instance, err := abi.NewToken(common.HexToAddress(smartContractAddress), ethClient)
	if err != nil {
		fmt.Println(err)
		return
	}

	localBalance, err := instance.BalanceOf(nil, common.HexToAddress(address))
	if err != nil {
		fmt.Println(err)
		return
	}

	r.Set(smartContractAddress, localBalance.String())
}
