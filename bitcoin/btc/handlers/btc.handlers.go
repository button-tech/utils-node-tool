package handlers

import (
	"log"
	"net/http"
	"strconv"

	"sync"

	"fmt"
	"github.com/button-tech/utils-node-tool/shared/db"
	"github.com/button-tech/utils-node-tool/shared/multiBalance"
	"github.com/button-tech/utils-node-tool/shared/responseModels"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"os"
)

func GetBalance(c *gin.Context) {

	address := c.Param("address")

	type BTC struct {
		Balance string `json:"balance"`
	}

	var btc BTC

	response := new(responses.BalanceResponse)

	balance, err := req.Get(os.Getenv("btc-api") + "/v1/address/" + address)
	if err != nil || balance.Response().StatusCode != 200 {
		endPoint, err := db.GetEndpoint("btc")
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		balance, err = req.Get(endPoint + "/addr/" + address + "/balance")
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		balanceFloat, err := strconv.ParseFloat(balance.String(), 64)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		balanceFloat *= 0.00000001

		balanceStr := fmt.Sprintf("%f", balanceFloat)

		response.Balance = balanceStr

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

func GetTxFee(c *gin.Context) {

	resp := new(responses.TransactionFeeResponse)

	// (148 * 1(input) + 34 * 2 (output))/1000 * 0.0001(minimal BTC)
	resp.Fee = 0.218 * 0.0001

	c.JSON(http.StatusOK, resp)
}

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
