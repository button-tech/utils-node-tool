package handlers

import (
	"context"
	"encoding/json"
	"github.com/button-tech/utils-node-tool/logger"
	"github.com/button-tech/utils-node-tool/nodetools"
	b "github.com/button-tech/utils-node-tool/nodetools"
	"github.com/button-tech/utils-node-tool/nodetools/storage"
	"github.com/button-tech/utils-node-tool/types/requests"
	"github.com/button-tech/utils-node-tool/types/responses"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/onrik/ethrpc"
	"github.com/qiangxue/fasthttp-routing"
	"math"
	"os"
	"time"
)

func GetBalance(c *routing.Context) error {

	start := time.Now()

	address := c.Param("address")

	balance, err := b.GetEtherBalance(address)
	if err != nil {
		logger.HandlerError("GetBalance", err)
		return err
	}

	response := new(responses.BalanceResponse)

	response.Balance = balance

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	logger.LogRequest(time.Since(start), os.Getenv("BLOCKCHAIN"), "GetBalance")

	return nil
}

func GetTxFee(c *routing.Context) error {

	start := time.Now()

	ethClient := ethrpc.New(storage.EndpointForReq.Get())

	gasPrice, err := ethClient.EthGasPrice()
	if err != nil {
		logger.HandlerError("GetTxFee", err)
		return err
	}

	fee := float64(gasPrice.Int64()*21000) / math.Pow(10, 18)

	response := new(responses.TransactionFeeResponse)

	response.Fee = fee

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	logger.LogRequest(time.Since(start), os.Getenv("BLOCKCHAIN"), "GetTxFee")

	return nil
}

func GetGasPrice(c *routing.Context) error {

	start := time.Now()

	ethClient := ethrpc.New(storage.EndpointForReq.Get())

	gasPrice, err := ethClient.EthGasPrice()
	if err != nil {
		logger.HandlerError("GetGasPrice", err)
		return err
	}

	response := new(responses.GasPriceResponse)

	response.GasPrice = gasPrice.Int64()

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	logger.LogRequest(time.Since(start), os.Getenv("BLOCKCHAIN"), "GetGasPrice")

	return nil
}

func GetTokenBalance(c *routing.Context) error {

	start := time.Now()

	userAddress := c.Param("user-address")

	smartContractAddress := c.Param("smart-contract-address")

	balance, err := b.GetTokenBalance(userAddress, smartContractAddress)
	if err != nil {
		logger.HandlerError("GetTokenBalance", err)
		return err
	}

	response := new(responses.BalanceResponse)

	response.Balance = balance

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	logger.LogRequest(time.Since(start), os.Getenv("BLOCKCHAIN"), "GetTokenBalance")

	return nil
}

func GetEstimateGas(c *routing.Context) error {

	start := time.Now()

	var data requests.EthEstimateGasRequest

	if err := json.Unmarshal(c.PostBody(), &data); err != nil {
		return err
	}

	gasLimit, err := nodetools.GetEstimateGas(&data)
	if err != nil {
		logger.HandlerError("GetEstimateGas", err)
		return err
	}

	response := new(responses.GasLimitResponse)

	response.GasLimit = gasLimit

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	logger.LogRequest(time.Since(start), os.Getenv("BLOCKCHAIN"), "GetEstimateGas")

	return nil
}

func GetNonce(c *routing.Context) error {

	start := time.Now()

	userAddress := c.Param("address")

	client, err := ethclient.Dial(storage.EndpointForReq.Get())
	if err != nil {
		return err
	}

	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(userAddress))
	if err != nil {
		logger.HandlerError("GetNonce", err)
		return err
	}

	result := struct {
		Nonce uint64 `json:"nonce"`
	}{
		Nonce: nonce,
	}

	if err := responses.JsonResponse(c, result); err != nil {
		return err
	}

	logger.LogRequest(time.Since(start), os.Getenv("BLOCKCHAIN"), "GetNonce")

	return nil
}
