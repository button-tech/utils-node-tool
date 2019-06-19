package handlers

import (
	"log"
	"math"
	"net/http"
	"sync"

	"github.com/button-tech/utils-node-tool/shared/db"
	"github.com/button-tech/utils-node-tool/shared/multiBalance"
	"github.com/button-tech/utils-node-tool/shared/responseModels"
	"github.com/gin-gonic/gin"
	"github.com/onrik/ethrpc"
	"os"
)

func GetBalance(c *gin.Context) {

	var etcClient = ethrpc.New(os.Getenv("etc-api"))

	balance, err := etcClient.EthGetBalance(c.Param("address"), "latest")

	if err != nil {

		reserveNode, err := db.GetReserveHost("etc")
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		etcClient = ethrpc.New(reserveNode)

		result, err := etcClient.EthGetBalance(c.Param("address"), "latest")
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		balance = result
	}

	response := new(responses.BalanceResponse)

	response.Balance = balance.String()

	c.JSON(http.StatusOK, response)
}

func GetTxFee(c *gin.Context) {

	etcClient := ethrpc.New(os.Getenv("etc-api"))

	gasPrice, err := etcClient.EthGasPrice()

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	fee := float64(gasPrice.Int64()*21000) / math.Pow(10, 18)

	response := new(responses.TransactionFeeResponse)

	response.Fee = fee

	c.JSON(http.StatusOK, response)
}

func GetGasPrice(c *gin.Context) {

	etcClient := ethrpc.New(os.Getenv("etc-api"))

	gasPrice, err := etcClient.EthGasPrice()

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	response := new(responses.GasPriceResponse)

	response.GasPrice = gasPrice.Int64()

	c.JSON(http.StatusOK, response)
}

func GetBalances(c *gin.Context) {

	type Request struct {
		AddressesArray []string `json:"addressesArray"`
	}

	req := new(Request)

	balances := multiBalance.New()

	err := c.BindJSON(&req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var wg sync.WaitGroup

	for i := 0; i < len(req.AddressesArray); i++ {
		wg.Add(1)
		go multiBalance.EtcWorker(&wg, req.AddressesArray[i], balances)
	}
	wg.Wait()

	response := new(responses.BalancesResponse)
	response.Balances = balances.Result

	c.JSON(http.StatusOK, response)
}
