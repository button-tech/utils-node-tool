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

func GetTxFee(c *gin.Context) {

	resp := new(responses.TransactionFeeResponse)

	// (148 * 1(input) + 34 * 2 (output))/1000 * 0.0001(minimal BCH)
	resp.Fee = 0.218 * 0.0001

	c.JSON(http.StatusOK, resp)
}

func GetUTXO(c *gin.Context) {

	address := c.Param("address")

	utxos, err := req.Get(os.Getenv("reserve-api") + "/v1/address/utxo/" + address)
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
