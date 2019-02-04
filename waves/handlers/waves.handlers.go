package handlers

import (
	"fmt"
	"github.com/button-tech/utils-node-tool/waves/responseModels"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"net/http"
	"os"
	"strconv"
)

var (
	wavesURL = os.Getenv("WAVES_NODE")
)

// @Summary Waves balance of account
// @Description return balance of account in Waves for specific node
// @Produce  application/json
// @Param   address        path    string     true        "address"
// @Success 200 {array} responses.BalanceResponse
// @Router /waves/balance/{address} [get]
// GetBalance return balance of account in Waves
func GetBalance(c *gin.Context) {

	type BalanceData struct {
		Balance int64 `json:"balance"`
	}

	address := c.Param("address")

	res, err := req.Get(wavesURL + "/addresses/balance/" + address)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	var data BalanceData

	res.ToJSON(&data)

	response := new(responses.BalanceResponse)

	response.Balance = strconv.FormatInt(data.Balance, 10)

	c.JSON(http.StatusOK, response)
}
