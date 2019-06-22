package shared

import (
	"context"
	"errors"
	"fmt"
	"github.com/button-tech/utils-node-tool/shared/abi"
	"github.com/button-tech/utils-node-tool/shared/db"
	"github.com/button-tech/utils-node-tool/shared/requests"
	"github.com/button-tech/utils-node-tool/shared/responses"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imroc/req"
	"github.com/onrik/ethrpc"
	"golang.org/x/crypto/sha3"
	"golang.org/x/sync/errgroup"
	"math/big"
	"os"
	"strconv"
	"sync"
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

	var ethClient = ethrpc.New(os.Getenv("main-api"))

	balance, err := ethClient.EthGetBalance(address, "latest")
	if err != nil {
		reserveNode, err := db.GetEndpoint("blockchain")
		if err != nil {
			return "", err
		}

		ethClient = ethrpc.New(reserveNode)

		result, err := ethClient.EthGetBalance(address, "latest")
		if err != nil {
			return "", err
		}

		balance = result
	}

	return balance.String(), nil
}

func GetTokenBalance(userAddress, smartContractAddress string) (string, error) {
	ethClient, err := ethclient.Dial(os.Getenv("main-api"))
	if err != nil {
		endPoint, err := db.GetEndpoint("blockchain")
		if err != nil {
			return "", err
		}

		ethClient, err = ethclient.Dial(endPoint)
		if err != nil {
			return "", err
		}
	}

	instance, err := abi.NewToken(common.HexToAddress(smartContractAddress), ethClient)
	if err != nil {
		return "", err
	}

	balance, err := instance.BalanceOf(nil, common.HexToAddress(userAddress))
	if err != nil {
		endPoint, err := db.GetEndpoint("main-api")
		if err != nil {
			return "", err
		}

		ethClient, err = ethclient.Dial(endPoint)
		if err != nil {
			return "", err
		}

		instance, err = abi.NewToken(common.HexToAddress(smartContractAddress), ethClient)
		if err != nil {
			return "", err
		}

		balance, err = instance.BalanceOf(nil, common.HexToAddress(userAddress))
		if err != nil {
			return "", err
		}
	}

	return balance.String(), nil
}

func GetEstimateGas(req *requests.EthEstimateGasRequest) (uint64, error) {

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

	ethClient, err := ethclient.Dial(os.Getenv("main-api"))
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

	var reserveUrl string

	currency := os.Getenv("blockchain")

	switch currency {
	case "btc":
		reserveUrl = "/addr/" + address + "/balance"
	case "bch":
		reserveUrl = "/v1/address/details/" + address
	case "ltc":
		reserveUrl = "/api/addr/" + address + "/balance"
	}

	data := struct {
		Balance string `json:"balance"`
	}{}

	reserveData := struct {
		Balance float64 `json:"balance"`
	}{}

	responseOfMainApi, err := req.Get(os.Getenv("main-api") + "/v1/address/" + address)

	if err != nil || responseOfMainApi.Response().StatusCode != 200 {
		endpoint, err := db.GetEndpoint(currency)
		if err != nil {
			return "", err
		}

		responseOfReserveApi, err := req.Get(endpoint + reserveUrl)
		if err != nil {
			return "", err
		}

		if responseOfReserveApi.Response().StatusCode != 200 {
			return "", errors.New("Bad request")
		}

		if currency == "bch" {
			err = responseOfReserveApi.ToJSON(&reserveData)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%f", reserveData.Balance), nil

		}

		balanceFloat, err := strconv.ParseFloat(responseOfReserveApi.String(), 64)
		if err != nil {
			return "", err
		}

		balanceFloat *= 0.00000001

		return fmt.Sprintf("%f", balanceFloat), nil
	}

	err = responseOfMainApi.ToJSON(&data)
	if err != nil {
		return "", err
	}

	return data.Balance, nil
}

func GetUtxo(address string) ([]responses.UTXO, error) {

	currency := os.Getenv("blockchain")

	var endPoint string

	var requestUrl string

	if currency == "bch" {
		endPoint = "https://rest.bitbox.earth"
	} else {
		nodeFromDB, err := db.GetEndpoint(currency)
		if err != nil {
			return nil, err
		}
		endPoint = nodeFromDB
	}

	switch currency {
	case "btc":
		requestUrl = endPoint + "/addr/" + address + "/utxo"
	case "bch":
		requestUrl = endPoint + "/v1/address/utxo/" + address
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
