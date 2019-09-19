package balance

import (
	"errors"
	"github.com/button-tech/utils-node-tool/shared/abi"
	"github.com/button-tech/utils-node-tool/utils-for-endpoints/estorage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imroc/req"
	"github.com/onrik/ethrpc"
	"golang.org/x/sync/errgroup"
	"os"
	"strconv"
	"sync"
	"time"
)

// Balances(by addresses list)
type Balances struct {
	sync.RWMutex
	AddressesAndBalances map[string]string
}

func NewBalances() *Balances {
	return &Balances{
		AddressesAndBalances: make(map[string]string),
	}
}

func (ds *Balances) SetBalances(address, balance string) {
	ds.Lock()
	ds.AddressesAndBalances[address] = balance
	ds.Unlock()
}

func GetUtxoBasedBalancesByList(addresses []string) (map[string]string, error) {

	result := NewBalances()

	var g errgroup.Group

	for _, address := range addresses {

		address := address

		g.Go(func() error {

			balance, err := GetUtxoBasedBalance(address)
			if err == nil {
				result.SetBalances(address, balance)
			}

			return err
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result.AddressesAndBalances, nil
}

// UTXO based blockchain - BTC, LTC, BCH
func GetUtxoBasedBalance(address string) (string, error) {

	s := struct {
		Balance interface{} `json:"balance"`
	}{}

	var result string

	res, err := req.Get(estorage.EndpointForReq.Get() + address)
	if err != nil || res.Response().StatusCode != 200 {
		result, err = UtxoBasedBalanceReq(address)
		if err != nil {
			return "", err
		}

		return result, nil
	}

	err = res.ToJSON(&s)
	if err != nil {
		return "", err
	}

	result, err = ParseUtxoApiResponse(s.Balance)
	if err != nil {
		return "", err
	}

	return result, nil
}

func UtxoBasedBalanceReq(address string) (string, error) {

	var endpoints []string

	currency := os.Getenv("BLOCKCHAIN")

	mainApi := os.Getenv("MAIN_API")

	mainUrl := mainApi + "/v1/address/" + address

	switch currency {
	case "btc":
		dbEndpoints := estorage.EndpointsFromDB.Get().Addresses
		for _, j := range dbEndpoints {
			j = j + "/addr/" + address
			endpoints = append(endpoints, j)
		}
		endpoints = append(endpoints, mainUrl)
	case "ltc":
		dbEndpoints := estorage.EndpointsFromDB.Get().Addresses
		for _, j := range dbEndpoints {
			j = j + "/api/addr/" + address
			endpoints = append(endpoints, j)
		}
		endpoints = append(endpoints, mainUrl)
	case "bch":
		dbEndpoints := estorage.EndpointsFromDB.Get().Addresses
		for _, j := range dbEndpoints {
			j = j + address
			endpoints = append(endpoints, j)
		}
		endpoints = append(endpoints, mainUrl)
	}

	balanceChan := make(chan string, len(endpoints))

	for _, addr := range endpoints {
		go func(addr string) {
			s := struct {
				Balance interface{} `json:"balance"`
			}{}

			res, err := req.Get(addr)
			if err != nil || res.Response().StatusCode != 200 {
				return
			}

			err = res.ToJSON(&s)
			if err != nil {
				return
			}

			balance, err := ParseUtxoApiResponse(s.Balance)
			if err != nil {
				return
			}

			balanceChan <- balance
		}(addr)
	}

	select {
	case result := <-balanceChan:
		return result, nil
	case <-time.After(2 * time.Second):
		return "", errors.New("Bad request")
	}
}

// ETH based
func GetEtherBalance(address string) (string, error) {

	ethClient := ethrpc.New(estorage.EndpointForReq.Get())

	res, err := ethClient.EthGetBalance(address, "latest")
	if err != nil {
		balance, err := EtherBalanceReq(address)
		if err != nil {
			return "", err
		}

		return balance, nil
	}

	return res.String(), nil

}

func GetTokenBalance(userAddress, smartContractAddress string) (string, error) {

	ethClient, err := ethclient.Dial(estorage.EndpointForReq.Get())
	if err != nil {
		return "", err
	}

	instance, err := abi.NewToken(common.HexToAddress(smartContractAddress), ethClient)
	if err != nil {
		return "", err
	}

	res, err := instance.BalanceOf(nil, common.HexToAddress(userAddress))
	if err != nil {
		balance, err := TokenBalanceReq(userAddress, smartContractAddress)
		if err != nil {
			return "", err
		}

		return balance, nil
	}

	return res.String(), nil
}

func TokenBalanceReq(userAddress, smartContractAddress string) (string, error) {

	endpoints := estorage.EndpointsFromDB.Get().Addresses

	endpoints = append(endpoints, os.Getenv("MAIN_API"))

	balanceChan := make(chan string, len(endpoints))

	for _, e := range endpoints {
		go func(e string) {
			ethClient, err := ethclient.Dial(e)
			if err != nil {
				return
			}

			instance, err := abi.NewToken(common.HexToAddress(smartContractAddress), ethClient)
			if err != nil {
				return
			}

			res, err := instance.BalanceOf(nil, common.HexToAddress(userAddress))
			if err != nil {
				return
			}

			balanceChan <- res.String()

		}(e)
	}

	select {
	case result := <-balanceChan:
		return result, nil
	case <-time.After(2 * time.Second):
		return "", errors.New("Bad request")
	}
}

func EtherBalanceReq(address string) (string, error) {

	endpoints := estorage.EndpointsFromDB.Get().Addresses

	endpoints = append(endpoints, os.Getenv("MAIN_API"))

	balanceChan := make(chan string, len(endpoints))

	for _, e := range endpoints {
		go func(e string) {
			ethClient := ethrpc.New(e)
			res, err := ethClient.EthGetBalance(address, "latest")
			if err != nil {
				return
			}

			balanceChan <- res.String()
		}(e)
	}

	select {
	case result := <-balanceChan:
		return result, nil
	case <-time.After(2 * time.Second):
		return "", errors.New("Bad request")
	}
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
