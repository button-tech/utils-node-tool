package handlers

import (
	"fmt"
	"github.com/button-tech/utils-node-tool/btc/handlers/responseModels"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"net/http"
	"os"
)

var (
	btcURL = os.Getenv("BTC_NODE")
)

// @Summary BTC balance of account
// @Description return balance of account in BTC for specific node
// @Produce  application/json
// @Param   address        path    string     true        "address"
// @Router /btc/balance/{address} [get]
// GetBalance return balance of account in BTC for specific node
func GetBalance(c *gin.Context) {

	address := c.Param("address")

	utxos, err := req.Get(btcURL + "/api/addr/" + address + "/utxo")
	if err != nil {
		fmt.Println(err)
	}

	type UTXO struct {
		Address string `json:"address"`
		Txid string `json:txid"`
		Vout int `json:"vout"`
		ScriptPubKey string `json:"scriptPubKey"`
		Amount float64 `json:"amount"`
		Satoshis int `json:"satoshis"`
		Height int `json:"height"`
		Confirmations int `json:"confirmations"`
	}

	var respArr []UTXO
	utxos.ToJSON(&respArr)

	var balance float64

	for _, j := range respArr {
		balance += j.Amount
	}

	response := new(responses.BalanceResponse)

	response.Balance = balance

	c.JSON(http.StatusOK, response)

}
