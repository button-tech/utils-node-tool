package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/button-tech/utils-node-tool/btc/handlers/responseModels"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
)

var (
	btcURL = os.Getenv("BTC_NODE")
)

// @Summary BTC balance of account
// @Description return balance of account in BTC for specific node
// @Produce  application/json
// @Param   address        path    string     true        "address"
// @Success 200 {array} responses.BalanceResponse
// @Router /btc/balance/{address} [get]
// GetBalance return balance of account in BTC for specific node
func GetBalance(c *gin.Context) {

	address := c.Param("address")

	balance, err := req.Get(btcURL + "/insight-api/addr/" + address + "/balance")
	if err != nil {
		fmt.Println(err)
	}

	balanceFloat, _ := strconv.ParseFloat(balance.String(), 64)

	balanceFloat *= 0.00000001

	response := new(responses.BalanceResponse)

	response.Balance = balanceFloat

	c.JSON(http.StatusOK, response)

}

// @Summary BTC UTXO of account
// @Description return UTXO of account
// @Produce  application/json
// @Param   address        path    string     true        "address"
// @Success 200 {array} responses.UTXOResponse
// @Router /btc/utxo/{address} [get]
// GetUTXO return UTXO of account
func GetUTXO(c *gin.Context) {

	address := c.Param("address")
	utxos, err := req.Get(btcURL + "/insight-api/addr/" + address + "/utxo")
	if err != nil {
		fmt.Println(err)
	}

	var respArr []responses.UTXO

	utxos.ToJSON(&respArr)

	response := new(responses.UTXOResponse)

	response.Utxo = respArr

	c.JSON(http.StatusOK, response)
}
