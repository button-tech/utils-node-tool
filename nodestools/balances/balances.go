package balances

import (
	"errors"
	"github.com/button-tech/utils-node-tool/abi"
	"github.com/button-tech/utils-node-tool/requests"
	"github.com/button-tech/utils-node-tool/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imroc/req"
	"github.com/onrik/ethrpc"
	"golang.org/x/sync/errgroup"
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

	var (
		s      requests.UtxoBasedBalance
		result string
	)

	res, err := req.Get(storage.EndpointForReq.Get() + "/address/" + address)
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

	return s.Balance, nil
}

func UtxoBasedBalanceReq(address string) (string, error) {

	endpoints := storage.EndpointsFromDB.Get().Addresses

	balanceChan := make(chan string, len(endpoints))

	for _, addr := range endpoints {
		go func(addr string) {

			var s requests.UtxoBasedBalance

			res, err := req.Get(addr + "/address/" + address)
			if err != nil || res.Response().StatusCode != 200 {
				return
			}

			err = res.ToJSON(&s)
			if err != nil {
				return
			}

			balanceChan <- s.Balance

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

	ethClient := ethrpc.New(storage.EndpointForReq.Get())

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

	ethClient, err := ethclient.Dial(storage.EndpointForReq.Get())
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

	endpoints := storage.EndpointsFromDB.Get().Addresses

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

	endpoints := storage.EndpointsFromDB.Get().Addresses

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
