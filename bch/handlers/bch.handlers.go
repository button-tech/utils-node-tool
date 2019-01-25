package handlers

import (
	"fmt"
	"github.com/button-tech/utils-node-tool/bch/handlers/responseModels"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"net/http"
	"os"
	"strconv"
)

var (
	bchURL = os.Getenv("BCH_NODE")
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

	balance, err := req.Get(bchURL + "/api/addr/" + address + "/balance")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	balanceFloat, _ := strconv.ParseFloat(balance.String(), 64)

	response := new(responses.BalanceResponse)

	response.Balance = balanceFloat

	c.JSON(http.StatusOK, response)

}

// @Summary BCH fee of tx
// @Description return fee of tx in BCH
// @Produce  application/json
// @Success 200 {array} responses.TransactionFeeResponse
// @Router /bch/transactionFee [get]
// GetBalance return fee of tx in BCH
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
	utxos, err := req.Get(bchURL + "/api/addr/" + address + "/utxo")
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
