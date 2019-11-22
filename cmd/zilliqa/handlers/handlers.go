package handlers

import (
	"errors"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/button-tech/utils-node-tool/logger"
	"github.com/button-tech/utils-node-tool/types/responses"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetBalance(c *gin.Context) {

	start := time.Now()

	zilliqaAddress := c.Param("address")

	decodedAddress, err := bech32.FromBech32Addr(zilliqaAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		logger.Error("GetBalance: bech32.FromBech32Addr", err.Error(), logger.Params{
			"address": zilliqaAddress,
		})
		return
	}

	endpoint := provider.NewProvider("https://api.zilliqa.com/")
	if endpoint == nil {
		c.JSON(http.StatusInternalServerError, errors.New("api.zilliqa.com isn't available now"))
		logger.Error("GetBalance: provider.NewProvider", "endpoint==nil", logger.Params{
			"api": "https://api.zilliqa.com",
		})
		return
	}

	balance := endpoint.GetBalance(decodedAddress)
	if balance == nil {
		c.JSON(http.StatusInternalServerError, errors.New("Problems with api.zilliqa.com"))
		logger.Error("GetBalance: endpoint.GetBalance", "balance==nil", logger.Params{
			"address": decodedAddress,
		})
		return
	}

	response := new(responses.BalanceResponse)

	if balance.Result == nil {
		response.Balance = "0"
		c.JSON(http.StatusOK, response)
		logger.LogRequest(time.Since(start), "zilliqa", "GetBalance")
		return
	}

	response.Balance = fmt.Sprintf("%v", balance.Result.(map[string]interface{})["balance"])

	c.JSON(http.StatusOK, response)

	logger.LogRequest(time.Since(start), "zilliqa", "GetBalance")

}
