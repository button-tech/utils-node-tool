package handlers

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/button-tech/utils-node-tool/Shared/db"
	"github.com/button-tech/utils-node-tool/Shared/multiBalance"
	"github.com/button-tech/utils-node-tool/Shared/responseModels"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
)

// @Summary Waves balance of account
// @Description return balance of account in Waves for specific node
// @Produce  application/json
// @Param   address        path    string     true        "address"
// @Success 200 {array} responses.BalanceResponse
// @Router /waves/balance/{address} [get]
// GetBalance return balance of account in Waves
func GetBalance(c *gin.Context) {

	address := c.Param("address")

	endPoint, err := db.GetEndpoint("waves")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	res, err := req.Get(endPoint + "/addresses/balance/" + address)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var data responses.BalanceData

	err = res.ToJSON(&data)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	response := new(responses.BalanceResponse)

	response.Balance = strconv.FormatInt(data.Balance, 10)

	c.JSON(http.StatusOK, response)
}

// @Summary Waves balance of accounts by list
// @Description return balances of accounts in Waves
// @Produce  application/json
// @Param addressesArray     body string true "addressesArray"
// @Success 200 {array} responses.BalancesResponse
// @Router /waves/balances [post]
// GetBalanceForMultipleAdresses return balances of accounts in Waves
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
		go multiBalance.WavesWorker(&wg, request.AddressesArray[i], balances)
	}
	wg.Wait()

	response := new(responses.BalancesResponse)

	response.Balances = balances.Result

	c.JSON(http.StatusOK, response)
}

// @Summary Waves fee
// @Description return Waves fee
// @Produce  application/json
// @Success 200 {array} responses.TransactionFeeResponse
// @Router /waves/transactionFee [get]
// GetBalance return Waves fee
func GetTxFee(c *gin.Context) {

	resp := new(responses.TransactionFeeResponse)

	resp.Fee = 0.001

	c.JSON(http.StatusOK, resp)
}
