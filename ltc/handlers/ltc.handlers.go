package handlers

import (
	"fmt"
	"github.com/button-tech/utils-node-tool/ltc/handlers/multi-balance"
	"github.com/button-tech/utils-node-tool/ltc/handlers/responseModels"
	"github.com/button-tech/utils-node-tool/ltc/handlers/storage"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"net/http"
	"sync"
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

	balance, err := req.Get(storage.LtcURL + "/api/addr/" + address + "/balance")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	response := new(responses.BalanceResponse)

	response.Balance = balance.String()

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
	utxos, err := req.Get(storage.LtcURL + "/api/addr/" + address + "/utxo")
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
