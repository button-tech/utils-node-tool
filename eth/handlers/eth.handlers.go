package handlers

import (
	"log"
	"math"
	"os"
	"encoding/json"
	"github.com/button-tech/utils-node-tool/shared/responseModels"
	"github.com/onrik/ethrpc"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/button-tech/utils-node-tool/eth/ethUtils"
)

func GetBalance(c *routing.Context) error {

	address := c.Param("address")

	balance, err := ethUtils.GetBalance(address)
	if err != nil{
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

	balance, err := ethUtils.GetTokenBalance(userAddress, smartContractAddress)
	if err != nil{
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

	var txData ethUtils.TxData

	if err := json.Unmarshal(c.PostBody(), &txData); err != nil {
		log.Println(err)
		return err
	}

	gasLimit, err := ethUtils.GetEstimateGas(&txData)
	if err != nil{
		return err
	}

	response := new(responses.GasLimitResponse)

	response.GasLimit = gasLimit

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	return nil
}

//func GetBalances(c *gin.Context) {
//
//	type Request struct {
//		AddressesArray []string `json:"addressesArray"`
//	}
//
//	req := new(Request)
//
//	balances := multiBalance.New()
//
//	err := c.BindJSON(&req)
//	if err != nil {
//		log.Println(err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
//		return
//	}
//
//	var wg sync.WaitGroup
//
//	for i := 0; i < len(req.AddressesArray); i++ {
//		wg.Add(1)
//		go multiBalance.EthWorker(&wg, req.AddressesArray[i], balances)
//	}
//	wg.Wait()
//
//	response := new(responses.BalancesResponse)
//	response.Balances = balances.Result
//
//	c.JSON(http.StatusOK, response)
//}
//
//func GetTokenBalances(c *gin.Context) {
//
//	type Request struct {
//		OwnerAddress   string   `json:"ownerAddress"`
//		SmartAddresses []string `json:"smartAddresses"`
//	}
//
//	req := new(Request)
//
//	balances := multiBalance.New()
//
//	err := c.BindJSON(&req)
//	if err != nil {
//		log.Println(err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
//		return
//	}
//
//	var wg sync.WaitGroup
//
//	for i := 0; i < len(req.SmartAddresses); i++ {
//		wg.Add(1)
//		go multiBalance.TokenWorker(&wg, req.OwnerAddress, req.SmartAddresses[i], balances)
//	}
//	wg.Wait()
//
//	response := new(responses.BalancesResponse)
//	response.Balances = balances.Result
//
//	c.JSON(http.StatusOK, response)
//}
