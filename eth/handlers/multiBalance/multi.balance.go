package multiBalance

import (
	"fmt"
	"github.com/button-tech/utils-node-tool/eth/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"sync"
	"github.com/button-tech/utils-node-tool/db"
	"log"
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

	endPoint, err := db.GetEndpoint("eth")
	if err != nil{
		log.Println(err)
		return
	}

	ethClient := ethrpc.New(endPoint)

	balance, err := ethClient.EthGetBalance(addr, "latest")
	if err != nil {
		fmt.Println(err)
		return
	}

	r.Set(addr, balance.String())
}

func TokenWorker(wg *sync.WaitGroup, address string, smartContractAddress string, r *Data) {

	defer wg.Done()

	endPoint, err := db.GetEndpoint("eth")
	if err != nil{
		log.Println(err)
		return
	}

	ethClient, err := ethclient.Dial(endPoint)
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
