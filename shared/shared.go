package shared

import (
	"context"
	"errors"
	"math/big"
	"os"
	"sync"

	"github.com/button-tech/utils-node-tool/shared/abi"
	"github.com/button-tech/utils-node-tool/shared/db"
	"github.com/button-tech/utils-node-tool/shared/requests"
	"github.com/button-tech/utils-node-tool/shared/responses"
	"github.com/button-tech/utils-node-tool/utils-for-endpoints/estorage"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imroc/req"
	"github.com/onrik/ethrpc"
	"golang.org/x/crypto/sha3"
	"golang.org/x/sync/errgroup"
	"log"
	"strconv"
	"strings"
	"sync/atomic"
)

// Balances(by addresses list)

type Balances struct {
	sync.Mutex
	AddressesAndBalances map[string]string
}

func NewBalancesMap() *Balances {
	return &Balances{
		AddressesAndBalances: make(map[string]string),
	}
}

func (ds *Balances) Set(address, balance string) {
	ds.Lock()
	ds.AddressesAndBalances[address] = balance
	ds.Unlock()
}

func GetUtxoBasedBalancesByList(addresses []string) (map[string]string, error) {

	result := NewBalancesMap()

	var g errgroup.Group

	for _, address := range addresses {

		address := address

		g.Go(func() error {

			balance, err := GetUtxoBasedBalance(address)
			if err == nil {
				result.Set(address, balance)
			}

			return err
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result.AddressesAndBalances, nil
}

// ETH based
func GetEthBasedBalance(address string) (string, error) {

	currency := os.Getenv("blockchain")

	var counter int32

	endpoints := estorage.EndpointsFromDB.Get(currency).Addresses
	endpoints = append(endpoints, os.Getenv("main-api"))

	var balance string
	result := make(chan string)

	for _, addr := range endpoints {
		addr := addr
		go func() {
			ethClient := ethrpc.New(addr)
			res, err := ethClient.EthGetBalance(address, "latest")
			if err != nil {
				atomic.AddInt32(&counter, 1)
				return
			}

			result <- res.String()
		}()
	}

LOOP:
	for {
		select {
		case balance = <-result:
			break LOOP
		default:
			if int(counter) == len(endpoints) {
				return "", errors.New("Bad request")
			}
		}
	}

	return balance, nil
}

func GetTokenBalance(userAddress, smartContractAddress string) (string, error) {
	currency := os.Getenv("blockchain")

	var counter int32

	endpoints := estorage.EndpointsFromDB.Get(currency).Addresses
	endpoints = append(endpoints, os.Getenv("main-api"))

	var balance string
	result := make(chan string)

	for _, addr := range endpoints {
		addr := addr
		go func() {
			ethClient, err := ethclient.Dial(addr)
			if err != nil {
				atomic.AddInt32(&counter, 1)
				return
			}

			instance, err := abi.NewToken(common.HexToAddress(smartContractAddress), ethClient)
			if err != nil {
				atomic.AddInt32(&counter, 1)
				return
			}

			res, err := instance.BalanceOf(nil, common.HexToAddress(userAddress))
			if err != nil {
				atomic.AddInt32(&counter, 1)
				return
			}

			result <- res.String()
		}()
	}

LOOP:
	for {
		select {
		case balance = <-result:
			break LOOP
		default:
			if int(counter) == len(endpoints) {
				return "", errors.New("Bad request")
			}
		}
	}

	return balance, nil
}

func GetEstimateGas(req *requests.EthEstimateGasRequest) (uint64, error) {

	currency := os.Getenv("blockchain")

	toAddress := common.HexToAddress(req.ToAddress)

	tokenAddress := common.HexToAddress(req.TokenAddress)

	amount := new(big.Int)
	amount.SetString(req.Amount, 10)

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()

	_, err := hash.Write(transferFnSignature)
	if err != nil {
		return 0, err
	}

	methodID := hash.Sum(nil)[:4]

	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	endPoint, err := estorage.GetEndpoint(currency)
	if err != nil {
		return 0, err
	}

	ethClient, err := ethclient.Dial(endPoint)
	if err != nil {
		return 0, err
	}

	gasLimit, err := ethClient.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	})

	if err != nil {
		return 0, err
	}

	return gasLimit, nil
}

// UTXO based blockchain - BTC, LTC, BCH
func GetUtxoBasedBalance(address string) (string, error) {

	var endpoints []string

	currency := os.Getenv("blockchain")

	mainApi := os.Getenv("main-api")

	mainUrl := mainApi + "/v1/address/" + address

	switch currency {
	case "btc":
		dbEndpoints := estorage.EndpointsFromDB.BtcEndpoints.Addresses
		for _, j := range dbEndpoints {
			j = j + "/addr/" + address
			endpoints = append(endpoints, j)
		}
		endpoints = append(endpoints, mainUrl)
	case "ltc":
		dbEndpoints := estorage.EndpointsFromDB.LtcEndpoints.Addresses
		for _, j := range dbEndpoints {
			j = j + "/api/addr/" + address
			endpoints = append(endpoints, j)
		}
		endpoints = append(endpoints, mainUrl)
	case "bch":
		endpoints = append(endpoints, mainUrl)
		endpoints = append(endpoints, "https://rest.bitbox.earth/v1/address/details/"+address)
	}

	balance, err := FastUxoBasedReq(endpoints)
	if err != nil {
		return "", err
	}

	return balance, nil
}

func FastUxoBasedReq(endpoints []string) (string, error) {
	result := make(chan string)

	var counter int32

	var balance string

	for _, addr := range endpoints {
		addr := addr
		go func() {
			s := struct {
				Balance interface{} `json:"balance"`
			}{}
			res, err := req.Get(addr)
			if err != nil || res.Response().StatusCode != 200 {
				atomic.AddInt32(&counter, 1)
				return
			}

			err = res.ToJSON(&s)
			if err != nil {
				atomic.AddInt32(&counter, 1)
				return
			}

			balance, err := ParseUtxoApiResponse(s.Balance)
			if err != nil {
				atomic.AddInt32(&counter, 1)
				return
			}
			result <- balance
		}()
	}

LOOP:
	for {
		select {
		case balance = <-result:
			break LOOP
		default:
			if int(counter) == len(endpoints) {
				return "", errors.New("Bad request")
			}
		}
	}

	return balance, nil
}

func GetUtxo(address string) ([]responses.UTXO, error) {

	currency := os.Getenv("blockchain")

	var requestUrl, endPoint string

	var err error

	if currency != "bch"{
		endPoint, err = estorage.GetEndpoint(currency)
		if err != nil {
			return nil, err
		}
	}

	switch currency {
	case "btc":
		requestUrl = endPoint + "/addr/" + address + "/utxo"
	case "bch":
		requestUrl = "https://rest.bitbox.earth/v1/address/utxo/" + address
	case "ltc":
		requestUrl = endPoint + "/api/addr/" + address + "/utxo"
	}

	utxos, err := req.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	if utxos.Response().StatusCode != 200 {
		return nil, errors.New("Bad request")
	}

	var utxoArray []responses.UTXO

	err = utxos.ToJSON(&utxoArray)
	if err != nil {
		return nil, err
	}

	return utxoArray, nil
}

func GetEthBasedBlockNumber(currency, addr string) (int64, error) {
	header := req.Header{
		"Content-Type": "application/json",
	}

	params := strings.NewReader("{\n\"jsonrpc\":\"2.0\",\n\"method\":\"eth_getBlockByNumber\",\n\"params\":[\"latest\", false],\n\"id\":1\n}")

	resp, err := req.Post(addr, header, params)

	if err != nil || resp.Response().StatusCode != 200 {
		err := DeleteEntry(currency, addr)
		if err != nil {
			return 0, err
		}
		return 0, nil
	}

	info := struct {
		Result struct {
			Number string `json:"number"`
		}
	}{}

	err = resp.ToJSON(&info)
	if err != nil {
		return 0, err
	}

	if len(info.Result.Number) == 0 {
		return 0, errors.New("Bad request")
	}

	hexNumber := []byte(info.Result.Number)

	intNumber, err := strconv.ParseInt(string(hexNumber[2:]), 16, 64)
	if err != nil {
		return 0, err
	}

	return intNumber, nil
}

func GetUtxoBasedBlockNumber(currency, addr string) (int64, error) {
	var url string

	switch currency {
	case "btc":
		url = "/sync"
	case "ltc":
		url = "/api/sync"
	}

	res, err := req.Get(addr + url)
	if err != nil || res.Response().StatusCode != 200 {
		err := DeleteEntry(currency, addr)
		if err != nil {
			return 0, err
		}
		log.Println("Status code:" + strconv.Itoa(res.Response().StatusCode))
	}

	info := struct {
		BlockChainHeight int64 `json:"blockChainHeight"`
	}{}

	err = res.ToJSON(&info)
	if err != nil {
		return 0, err
	}

	return info.BlockChainHeight, nil
}

func Max(array []int64) int64 {
	var max int64 = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
	}
	return max
}

func ParseUtxoApiResponse(i interface{}) (string, error) {
	switch i.(type) {
	case string:
		return i.(string), nil
	case float64:
		return strconv.FormatFloat(i.(float64), 'f', 8, 64), nil
	}
	return "", errors.New("Bad request")
}

func DeleteEntry(currency, address string) error {
	isDel, err := db.AddToStoppedList(currency, address)
	if err != nil {
		return err
	}
	if !isDel {
		return errors.New("Can't del!\n")
	} else {
		log.Printf("Add to stopped list %s: %s", currency, address)
	}
	return nil
}
