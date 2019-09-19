package handlers

import (
	"encoding/json"
	"github.com/button-tech/utils-node-tool/shared"
	b "github.com/button-tech/utils-node-tool/shared/balance"
	"github.com/button-tech/utils-node-tool/shared/requests"
	"github.com/button-tech/utils-node-tool/shared/responses"
	"github.com/button-tech/utils-node-tool/utils-for-endpoints/estorage"
	"github.com/onrik/ethrpc"
	"github.com/qiangxue/fasthttp-routing"
	"log"
	"math"
	"os"
	"runtime"
)

func GetBalance(c *routing.Context) error {

	address := c.Param("address")

	balance, err := b.GetEtherBalance(address)
	if err != nil {
		return err
	}

	response := new(responses.BalanceResponse)

	response.Balance = balance

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	log.Println(runtime.NumGoroutine())

	return nil
}

func GetTxFee(c *routing.Context) error {

	ethClient := ethrpc.New(os.Getenv("MAIN_API"))

	gasPrice, err := ethClient.EthGasPrice()

	if err != nil {
		log.Println(err)
		return err
	}

	fee := float64(gasPrice.Int64()*21000) / math.Pow(10, 18)

	response := new(responses.TransactionFeeResponse)

	response.Fee = fee

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	return nil
}

func GetGasPrice(c *routing.Context) error {

	var ethClient *ethrpc.EthRPC

	switch os.Getenv("BLOCKCHAIN") {
	case "eth":
		ethClient = ethrpc.New(os.Getenv("MAIN_API"))
	case "etc":
		endPoint, err := estorage.GetEndpoint("etc")
		if err != nil {
			return err
		}
		ethClient = ethrpc.New(endPoint)
	}

	gasPrice, err := ethClient.EthGasPrice()
	if err != nil {
		log.Println(err)
		return err
	}

	response := new(responses.GasPriceResponse)

	response.GasPrice = gasPrice.Int64()

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	return nil
}

func GetTokenBalance(c *routing.Context) error {

	userAddress := c.Param("user-address")

	smartContractAddress := c.Param("smart-contract-address")

	balance, err := b.GetTokenBalance(userAddress, smartContractAddress)
	if err != nil {
		return err
	}

	response := new(responses.BalanceResponse)

	response.Balance = balance

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	return nil
}

func GetEstimateGas(c *routing.Context) error {

	var data requests.EthEstimateGasRequest

	if err := json.Unmarshal(c.PostBody(), &data); err != nil {
		log.Println(err)
		return err
	}

	gasLimit, err := shared.GetEstimateGas(&data)
	if err != nil {
		return err
	}

	response := new(responses.GasLimitResponse)

	response.GasLimit = gasLimit

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	return nil
}
