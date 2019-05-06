package multiBalance

import (
	"github.com/button-tech/utils-node-tool/shared/abi"
	"github.com/button-tech/utils-node-tool/shared/db"
	"github.com/button-tech/utils-node-tool/shared/responseModels"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imroc/req"
	"github.com/onrik/ethrpc"
	"log"
	"strconv"
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

func BtcWorker(wg *sync.WaitGroup, addr string, r *Data) {
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

func BchWorker(wg *sync.WaitGroup, addr string, r *Data) {
	defer wg.Done()

	endPoint, err := db.GetEndpoint("bch")
	if err != nil {
		log.Println(err)
		return
	}

	balance, err := req.Get(endPoint + "/api/addr/" + addr + "/balance")
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

func EtcWorker(wg *sync.WaitGroup, addr string, r *Data) {
	defer wg.Done()

	endPoint, err := db.GetEndpoint("etc")
	if err != nil {
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

func EthWorker(wg *sync.WaitGroup, addr string, r *Data) {
	defer wg.Done()

	endPoint, err := db.GetEndpoint("eth")
	if err != nil {
		log.Println(err)
		return
	}

	ethClient := ethrpc.New(endPoint)

	balance, err := ethClient.EthGetBalance(addr, "latest")
	if err != nil {
		log.Println(err)
		return
	}

	r.Set(addr, balance.String())
}

func TokenWorker(wg *sync.WaitGroup, address string, smartContractAddress string, r *Data) {

	defer wg.Done()

	endPoint, err := db.GetEndpoint("eth")
	if err != nil {
		log.Println(err)
		return
	}

	ethClient, err := ethclient.Dial(endPoint)
	if err != nil {
		log.Println(err)
		return
	}

	instance, err := abi.NewToken(common.HexToAddress(smartContractAddress), ethClient)
	if err != nil {
		log.Println(err)
		return
	}

	localBalance, err := instance.BalanceOf(nil, common.HexToAddress(address))
	if err != nil {
		log.Println(err)
		return
	}

	r.Set(smartContractAddress, localBalance.String())
}

func LtcWorker(wg *sync.WaitGroup, addr string, r *Data) {
	defer wg.Done()

	endPoint, err := db.GetEndpoint("ltc")
	if err != nil {
		log.Println(err)
		return
	}

	balance, err := req.Get(endPoint + "/api/addr/" + addr + "/balance")
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

func WavesWorker(wg *sync.WaitGroup, addr string, r *Data) {
	defer wg.Done()

	endPoint, err := db.GetEndpoint("ltc")
	if err != nil {
		log.Println(err)
		return
	}

	balance, err := req.Get(endPoint + "/addresses/balance/" + addr)
	if err != nil {
		log.Println(err)
		return
	}

	var data responses.BalanceData

	err = balance.ToJSON(&data)
	if err != nil {
		log.Println(err)
	}

	balanceStr := strconv.FormatInt(data.Balance, 10)

	r.Set(addr, balanceStr)
}
