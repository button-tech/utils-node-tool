package handlers

import (
	"log"
	"net/http"
	"sync"

	"fmt"
	"github.com/button-tech/utils-node-tool/shared/multiBalance"
	"github.com/button-tech/utils-node-tool/shared/responseModels"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"os"
)

// @Summary BCH balance of account
// @Description return balance of account in BCH for specific node
// @Produce  application/json
// @Param   address        path    string     true        "address"
// @Success 200 {array} responses.BalanceResponse
// @Router /bch/balance/{address} [get]
// GetBalance return balance of account in BCH for specific node
func GetBalance(c *gin.Context) {

	address := c.Param("address")

	type BCH struct {
		Balance string `json:"balance"`
	}

	type reserveBCH struct {
		Balance float64 `json:"balance"`
	}

	var bch BCH
	var reserveBch reserveBCH

	response := new(responses.BalanceResponse)

	balance, err := req.Get(os.Getenv("bch-api") + "/v1/address/" + address)

	if err != nil || balance.Response().StatusCode != 200 {
		balance, err = req.Get(os.Getenv("reserve-api") + "/v1/address/details/" + address)
		if err != nil || balance.Response().StatusCode != 200 {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		err = balance.ToJSON(&reserveBch)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		result := fmt.Sprintf("%f", reserveBch.Balance)

		response.Balance = result

		c.JSON(http.StatusOK, response)

		return
	}

	err = balance.ToJSON(&bch)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	response.Balance = bch.Balance

	c.JSON(http.StatusOK, response)

}

// @Summary BCH fee
// @Description return BCH fee
// @Produce  application/json
// @Success 200 {array} responses.TransactionFeeResponse
// @Router /bch/transactionFee [get]
// GetBalance return BCH fee
func GetTxFee(c *gin.Context) {

	resp := new(responses.TransactionFeeResponse)

	// (148 * 1(input) + 34 * 2 (output))/1000 * 0.0001(minimal BCH)
	resp.Fee = 0.218 * 0.0001

	c.JSON(http.StatusOK, resp)
}

// @Summary BCH UTXO of account
// @Description return UTXO of account
// @Produce  application/json
// @Param   address        path    string     true        "address"
// @Success 200 {array} responses.UTXOResponse
// @Router /bch/utxo/{address} [get]
// GetUTXO return UTXO of account
func GetUTXO(c *gin.Context) {

	address := c.Param("address")

	utxos, err := req.Get(os.Getenv("reserve-api") + "/v1/address/utxo/" + address);
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

// @Description return balances of accounts in BCH
// @Produce  application/json
// @Param addressesArray     body string true "addressesArray"
// @Success 200 {array} responses.BalancesResponse
// @Router /bch/balances [post]
// GetBalanceForMultipleAdresses return balances of accounts in BCH
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
		go multiBalance.BchWorker(&wg, request.AddressesArray[i], balances)
	}
	wg.Wait()

	response := new(responses.BalancesResponse)

	response.Balances = balances.Result

	c.JSON(http.StatusOK, response)
}
