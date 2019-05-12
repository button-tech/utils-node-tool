package handlers

import (
	"log"
	"net/http"
	"strconv"

	"sync"

	"github.com/button-tech/utils-node-tool/shared/db"
	"github.com/button-tech/utils-node-tool/shared/multiBalance"
	"github.com/button-tech/utils-node-tool/shared/responseModels"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"os"
	"fmt"
)

// @Summary BTC balance of account
// @Description return balance of account in BTC for specific node
// @Produce  application/json
// @Param   address        path    string     true        "address"
// @Success 200 {array} responses.BalanceResponse
// @Router /btc/balance/{address} [get]
// GetBalance return balance of account in BTC for specific node
func GetBalance(c *gin.Context) {

	address := c.Param("address")

	type BTC struct {
		Balance string `json:"balance"`
	}

	var btc BTC

	response := new(responses.BalanceResponse)

	balance, err := req.Get(os.Getenv("btc-api") + "/v1/address/" + address)
	if err != nil {
		endPoint, err := db.GetEndpoint("btc")
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		balance, err = req.Get(endPoint + "/addr/" + address + "/balance")
		if err != nil{
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		balanceFloat, err := strconv.ParseFloat(balance.String(), 64)
		if err != nil{
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		balanceFloat *= 0.00000001

		balanceStr := fmt.Sprintf("%f", balanceFloat)

		response.Balance = balanceStr


		response.Balance = balance.String()

		c.JSON(http.StatusOK, response)

		return
	}

	err = balance.ToJSON(&btc)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	response.Balance = btc.Balance

	c.JSON(http.StatusOK, response)

}

// @Summary return Amount of BTC that you need to send a transaction
// @Description return Amount of BTC that you need to send a transaction
// @Produce  application/json
// @Success 200 {array} responses.TransactionFeeResponse
// @Router /btc/bestTransactionFee [get]
// GetBalance return Amount of BTC that you need to send a transaction
func GetBextTxFee(c *gin.Context) {

	type BTCFee struct {
		FastestFee  int `json:"fastestFee"`
		HalfHourFee int `json:"halfHourFee"`
		HourFee     int `json:"hourFee"`
	}

	var feeObj BTCFee

	fee, err := req.Get("https://bitcoinfees.earn.com/api/v1/fees/recommended")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	err = fee.ToJSON(&feeObj)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	feeFloat, err := strconv.ParseFloat(strconv.Itoa(feeObj.FastestFee), 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	response := new(responses.TransactionFeeResponse)
	response.Fee = feeFloat

	c.JSON(http.StatusOK, response)
}

// @Summary BTC fee
// @Description return BTC fee
// @Produce  application/json
// @Success 200 {array} responses.TransactionFeeResponse
// @Router /btc/transactionFee [get]
// GetBalance return BTC fee
func GetTxFee(c *gin.Context) {

	resp := new(responses.TransactionFeeResponse)

	// (148 * 1(input) + 34 * 2 (output))/1000 * 0.0001(minimal BTC)
	resp.Fee = 0.218 * 0.0001

	c.JSON(http.StatusOK, resp)
}

// @Summary BTC UTXO of account
// @Description return UTXO of account
// @Produce  application/json
// @Param   address        path    string     true        "address"
// @Success 200 {array} responses.UTXOResponse
// @Router /btc/utxo/{address} [get]
// GetUTXO return UTXO of account
func GetUTXO(c *gin.Context) {

	address := c.Param("address")

	endPoint, err := db.GetEndpoint("btc")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	utxos, err := req.Get(endPoint + "/addr/" + address + "/utxo")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var respArr []responses.UTXO

	err = utxos.ToJSON(&respArr)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	response := new(responses.UTXOResponse)

	response.Utxo = respArr

	c.JSON(http.StatusOK, response)
}

// @Summary BTC balance of accounts by list
// @Description return balances of accounts in BTC
// @Produce  application/json
// @Param addressesArray     body string true "addressesArray"
// @Success 200 {array} responses.BalancesResponse
// @Router /btc/balances [post]
// GetBalanceForMultipleAdresses return balances of accounts in BTC
func GetBalances(c *gin.Context) {

	type Request struct {
		AddressesArray []string `json:"addressesArray"`
	}

	request := new(Request)

	err := c.BindJSON(&request)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var wg sync.WaitGroup

	balances := multiBalance.New()

	for i := 0; i < len(request.AddressesArray); i++ {
		wg.Add(1)
		go multiBalance.BtcWorker(&wg, request.AddressesArray[i], balances)
	}
	wg.Wait()

	response := new(responses.BalancesResponse)
	response.Balances = balances.Result

	c.JSON(http.StatusOK, response)
}
