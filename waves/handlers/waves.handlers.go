package handlers

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/button-tech/utils-node-tool/shared/multiBalance"
	"github.com/button-tech/utils-node-tool/shared/responseModels"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
)

func GetBalance(c *gin.Context) {

	address := c.Param("address")

	res, err := req.Get("https://nodes.wavesplatform.com/addresses/balance/" + address)
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

func GetTxFee(c *gin.Context) {

	resp := new(responses.TransactionFeeResponse)

	resp.Fee = 0.001

	c.JSON(http.StatusOK, resp)
}
