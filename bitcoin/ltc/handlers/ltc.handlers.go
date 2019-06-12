package handlers

import (
	"log"
	"net/http"
	"sync"

	"fmt"
	"github.com/button-tech/utils-node-tool/shared/db"
	"github.com/button-tech/utils-node-tool/shared/multiBalance"
	"github.com/button-tech/utils-node-tool/shared/responseModels"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"os"
	"strconv"
)

// @Summary LTC balance of account
// @Description return balance of account in LTC for specific node
// @Produce  application/json
// @Param   address        path    string     true        "address"
// @Success 200 {array} responses.BalanceResponse
// @Router /ltc/balance/{address} [get]
// GetBalance return balance of account in LTC for specific node
func GetBalance(c *gin.Context) {

	address := c.Param("address")

	type LTC struct {
		Balance string `json:"balance"`
	}

	var ltc LTC

	response := new(responses.BalanceResponse)

	balance, err := req.Get(os.Getenv("ltc-api") + "/v1/address/" + address)
	if err != nil || balance.Response().StatusCode != 200 {
		endPoint, err := db.GetEndpoint("ltc")
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		balance, err := req.Get(endPoint + "/api/addr/" + address + "/balance")
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

	err = balance.ToJSON(&ltc)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	response.Balance = ltc.Balance

	c.JSON(http.StatusOK, response)
}

// @Summary LTC fee
// @Description return LTC fee
// @Produce  application/json
// @Success 200 {array} responses.TransactionFeeResponse
// @Router /ltc/transactionFee [get]
// GetBalance return LTC fee
func GetTxFee(c *gin.Context) {

	resp := new(responses.TransactionFeeResponse)

	// (148 * 1(input) + 34 * 2 (output))/1000 * 0.001(minimal LTC)
	resp.Fee = 0.218 * 0.002

	c.JSON(http.StatusOK, resp)
}

// @Summary LTC UTXO of account
// @Description return UTXO of account
// @Produce  application/json
// @Param   address        path    string     true        "address"
// @Success 200 {array} responses.UTXOResponse
// @Router /ltc/utxo/{address} [get]
// GetUTXO return UTXO of account
func GetUTXO(c *gin.Context) {

	address := c.Param("address")

	endPoint, err := db.GetEndpoint("ltc")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	utxos, err := req.Get(endPoint + "/api/addr/" + address + "/utxo")
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

// @Description return balances of accounts in LTC
// @Produce  application/json
// @Param addressesArray     body string true "addressesArray"
// @Success 200 {array} responses.BalancesResponse
// @Router /ltc/balances [post]
// GetBalanceForMultipleAdresses return balances of accounts in LTC
func GetBalances(c *gin.Context) {

	type Request struct {
		AddressesArray []string `json:"addressesArray"`
	}

	request := new(Request)

	balances := multiBalance.New()

	err := c.BindJSON(&request)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var wg sync.WaitGroup

	for i := 0; i < len(request.AddressesArray); i++ {
		wg.Add(1)
		go multiBalance.LtcWorker(&wg, request.AddressesArray[i], balances)
	}
	wg.Wait()

	response := new(responses.BalancesResponse)

	response.Balances = balances.Result

	c.JSON(http.StatusOK, response)
}
