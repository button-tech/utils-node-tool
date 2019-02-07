package multiBalance

import (
	"fmt"
	"github.com/button-tech/utils-node-tool/eth/abi"
	"github.com/button-tech/utils-node-tool/eth/handlers/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math"
	"strconv"
	"sync"
)

type Balances struct {
	sync.Mutex
	Result []map[string]float64 // eth -> ["account":balance] tokens -> ["smart contract addr": token balance]
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
	balance, err := storage.EthClient.EthGetBalance(addr, "latest")
	if err != nil {
		fmt.Println(err)
		return
	}
	floatBalance, _ := strconv.ParseFloat(balance.String(), 64)
	ethBalance := floatBalance / math.Pow(10, 18)
	r.Set(addr, ethBalance)
}

func TokenWorker(wg *sync.WaitGroup, address string, smartContractAddress string, r *Balances) {

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