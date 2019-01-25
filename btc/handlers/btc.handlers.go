package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/button-tech/utils-node-tool/btc/handlers/multiBalance"
	"github.com/button-tech/utils-node-tool/btc/handlers/responseModels"
	"github.com/button-tech/utils-node-tool/btc/handlers/storage"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"sync"
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

	balance, err := req.Get(storage.BtcURL + "/insight-api/addr/" + address + "/balance")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	balanceFloat, _ := strconv.ParseFloat(balance.String(), 64)

	balanceFloat *= 0.00000001

	response := new(responses.BalanceResponse)

	response.Balance = balanceFloat

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
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	fee.ToJSON(&feeObj)

	feeFloat, err := strconv.ParseFloat(strconv.Itoa(feeObj.FastestFee), 64)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
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
	utxos, err := req.Get(storage.BtcURL + "/insight-api/addr/" + address + "/utxo")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	var respArr []responses.UTXO

	utxos.ToJSON(&respArr)

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
func GetBalanceForMultipleAdresses(c *gin.Context) {

	type Request struct {
		AddressesArray []string `json:"addressesArray"`
	}

	req := new(Request)

	balances := multiBalance.New()

	c.BindJSON(&req)

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
