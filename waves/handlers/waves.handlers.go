package handlers

import (
	"fmt"
	"github.com/button-tech/utils-node-tool/waves/handlers/responseModels"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"net/http"
	"strconv"
	"github.com/button-tech/utils-node-tool/waves/handlers/storage"
	"sync"
	"github.com/button-tech/utils-node-tool/waves/handlers/multiBalance"
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

	res, err := req.Get(storage.WavesURL + "/addresses/balance/" + address)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	var data storage.BalanceData

	res.ToJSON(&data)

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
