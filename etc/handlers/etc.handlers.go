package handlers

import (
	"log"
	"math"
	"net/http"
	"sync"

	"github.com/button-tech/utils-node-tool/db"
	"github.com/button-tech/utils-node-tool/etc/handlers/multiBalance"
	"github.com/button-tech/utils-node-tool/etc/handlers/responseModels"
	"github.com/gin-gonic/gin"
	"github.com/onrik/ethrpc"
)

// @Summary ETC balance of account
// @Description return balance of account in ETC for specific node
// @Produce  application/json
// @Param   address        path    string     true        "address"
// @Success 200 {array} responses.BalanceResponse
// @Router /etc/balance/{address} [get]
// GetBalance return balance of account in ETC for specific node
func GetBalance(c *gin.Context) {

	endPoint, err := db.GetEndpoint("etc")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	var etcClient = ethrpc.New(endPoint)

	balance, err := etcClient.EthGetBalance(c.Param("address"), "latest")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	response := new(responses.BalanceResponse)

	response.Balance = balance.String()

	c.JSON(http.StatusOK, response)
}

// @Summary return Amount of ETC that you need to send a transaction
// @Description return Amount of ETC that you need to send a transaction
// @Produce  application/json
// @Success 200 {array} responses.TransactionFeeResponse
// @Router /etc/transactionFee [get]
// GetTxFee return Amount of ETC that you need to send a transaction
func GetTxFee(c *gin.Context) {

	endPoint, err := db.GetEndpoint("etc")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	etcClient := ethrpc.New(endPoint)

	gasPrice, err := etcClient.EthGasPrice()

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	fee := float64(gasPrice.Int64()*21000) / math.Pow(10, 18)

	response := new(responses.TransactionFeeResponse)
	response.Fee = fee

	c.JSON(http.StatusOK, response)
}

// @Summary return gas price of specific node
// @Description return Amount of ETC that you need to send a transaction
// @Produce  application/json
// @Success 200 {array} responses.GasPriceResponse
// @Router /etc/gasPrice [get]
// GetGasPrice return gas price of specific node
func GetGasPrice(c *gin.Context) {

	endPoint, err := db.GetEndpoint("etc")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	etcClient := ethrpc.New(endPoint)

	gasPrice, err := etcClient.EthGasPrice()

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	response := new(responses.GasPriceResponse)
	response.GasPrice = gasPrice.Int64()

	c.JSON(http.StatusOK, response)
}

// @Summary ETC balance of accounts by list
// @Description return balances of accounts in ETC
// @Produce  application/json
// @Param addressesArray     body string true "addressesArray"
// @Success 200 {array} responses.BalancesResponse
// @Router /etc/balances [post]
// GetBalanceForMultipleAdresses return balances of accounts in ETC
func GetBalances(c *gin.Context) {

	type Request struct {
		AddressesArray []string `json:"addressesArray"`
	}

	req := new(Request)

	balances := multiBalance.New()

	err := c.BindJSON(&req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	var wg sync.WaitGroup

	for i := 0; i < len(req.AddressesArray); i++ {
		wg.Add(1)
		go multiBalance.Worker(&wg, req.AddressesArray[i], balances)
	}
	wg.Wait()

	response := new(responses.BalancesResponse)
	response.Balances = balances.Result

	c.JSON(http.StatusOK, response)
}
