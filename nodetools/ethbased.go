package nodetools

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/button-tech/utils-node-tool/nodetools/abi"
	"github.com/button-tech/utils-node-tool/nodetools/storage"
	"github.com/button-tech/utils-node-tool/types/requests"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imroc/req"
	"github.com/onrik/ethrpc"
	"strconv"
	"strings"
	"time"
)

func GetEstimateGas(req *requests.EthEstimateGasRequest) (uint64, error) {

	ethClient, err := ethclient.Dial(storage.EndpointForReq.Get())
	if err != nil {
		return 0, err
	}

	address := common.HexToAddress(req.ContractAddress)

	data, err := hex.DecodeString(req.Data)
	if err != nil {
		return 0, err
	}

	gasLimit, err := ethClient.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &address,
		Data: data,
	})
	if err != nil {
		return 0, err
	}

	if gasLimit < 38000 {
		gasLimit = 80000
	}

	return gasLimit, nil
}

func GetEthBasedBlockNumber(addr string) (int64, error) {
	header := req.Header{
		"Content-Type": "application/json",
	}

	params := strings.NewReader("{\n\"jsonrpc\":\"2.0\",\n\"method\":\"eth_getBlockByNumber\",\n\"params\":[\"latest\", false],\n\"id\":1\n}")

	resp, err := req.Post(addr, header, params)
	if err != nil {
		return 0, err
	}

	if resp.Response().StatusCode != 200 {
		return 0, errors.New("Bad request")
	}

	var info requests.EthBasedBlocksHeight

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
