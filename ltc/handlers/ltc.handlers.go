package handlers

import (
	"fmt"
	"github.com/button-tech/utils-node-tool/ltc/handlers/responseModels"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"net/http"
	"os"
	"strconv"
)

var (
	ltcURL = os.Getenv("LTC_NODE")
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

	balance, err := req.Get(ltcURL + "/api/addr/" + address + "/balance")
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
	utxos, err := req.Get(ltcURL + "/api/addr/" + address + "/utxo")
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
