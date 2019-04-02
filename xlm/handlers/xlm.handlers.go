package handlers

import (
	"log"
	"net/http"

	"github.com/button-tech/utils-node-tool/xlm/handlers/responseModels"
	"github.com/button-tech/utils-node-tool/xlm/handlers/storage"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
)

// @Summary Stellar balance of account
// @Description return balance of account in Stellar for specific node
// @Produce  application/json
// @Param   address        path    string     true        "address"
// @Success 200 {array} responses.BalanceResponse
// @Router /stellar/balance/{address} [get]
// GetBalance return balance of account in Stellar for specific node
func GetBalance(c *gin.Context) {

	var balance storage.StellarBalance

	balanceReq, err := req.Get(storage.StellarNodeAddress.Address + "/accounts/" + c.Param("address"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	err = balanceReq.ToJSON(&balance)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	var stellarBalanceString string

	for _, j := range balance.Balances {
		if j.Asset_type == "native" {
			stellarBalanceString = j.Balance
		}
	}

	if stellarBalanceString == "" {
		stellarBalanceString = "0"
	}

	response := new(responses.BalanceResponse)

	response.Balance = stellarBalanceString

	c.JSON(http.StatusOK, response)

}
