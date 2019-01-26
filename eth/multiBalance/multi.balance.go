package multiBalance

import (
	"fmt"
	"github.com/button-tech/utils-node-tool/eth/abi"
	"github.com/button-tech/utils-node-tool/eth/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math"
	"strconv"
	"sync"
)

type Data struct {
	sync.Mutex
	Result map[string]float64
}

func New() *Data {
	return &Data{
		Result: make(map[string]float64),
	}
}

func (ds *Data) set(key string, value float64) {
	ds.Result[key] = value
}

func (ds *Data) Set(key string, value float64) {
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
	floatBalance, _ := strconv.ParseFloat(balance.String(), 64)
	ethBalance := floatBalance / math.Pow(10, 18)
	r.Set(addr, ethBalance)
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

	floatTokenBalance, _ := strconv.ParseFloat(localBalance.String(), 64)

	tokenBalance := floatTokenBalance / math.Pow(10, 18)

	r.Set(smartContractAddress, tokenBalance)
}
