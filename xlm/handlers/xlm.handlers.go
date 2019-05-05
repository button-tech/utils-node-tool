package handlers

import (
	"github.com/button-tech/utils-node-tool/db"
	"github.com/button-tech/utils-node-tool/xlm/handlers/responseModels"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"log"
	"net/http"
)

// @Summary Stellar balance of account
// @Description return balance of account in Stellar for specific node
// @Produce  application/json
// @Param   address        path    string     true        "address"
// @Success 200 {array} responses.BalanceResponse
// @Router /stellar/balance/{address} [get]
// GetBalance return balance of account in Stellar for specific node
func GetBalance(c *gin.Context) {

	type StellarBalance struct {
		Balances []struct {
			Balance             string `json:"balance"`
			Buying_liabilities  string `json:"buying_liabilities"`
			Selling_liabilities string `json:"selling_liabilities"`
			Asset_type          string `json:"asset_type"`
		}
	}

	var balance StellarBalance

	endPoint, err := db.GetEndpoint("xlm")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	balanceReq, err := req.Get(endPoint + "/accounts/" + c.Param("address"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	err = balanceReq.ToJSON(&balance)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
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
