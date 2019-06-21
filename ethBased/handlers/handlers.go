package handlers

import (
	"encoding/json"
	"github.com/button-tech/utils-node-tool/shared"
	"github.com/button-tech/utils-node-tool/shared/requests"
	"github.com/button-tech/utils-node-tool/shared/responses"
	"github.com/onrik/ethrpc"
	"github.com/qiangxue/fasthttp-routing"
	"log"
	"math"
	"os"
)

func GetBalance(c *routing.Context) error {

	address := c.Param("address")

	balance, err := shared.GetEthBasedBalance(address)
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

func GetTxFee(c *routing.Context) error {

	ethClient := ethrpc.New(os.Getenv("main-api"))

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

	ethClient := ethrpc.New(os.Getenv("main-api"))

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

	balance, err := shared.GetTokenBalance(userAddress, smartContractAddress)
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
