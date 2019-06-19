package handlers

import (
	"log"
	"net/http"
	"github.com/button-tech/utils-node-tool/shared/responseModels"
	"github.com/gin-gonic/gin"
	"github.com/button-tech/utils-node-tool/shared/btcUtils"
)

func GetBalance(c *gin.Context) {

	address := c.Param("address")

	response := new(responses.BalanceResponse)

	balance, err :=  btcUtils.GetBtcBlockChainBalance(address)
	if err != nil{
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	response.Balance = balance

	c.JSON(http.StatusOK, response)

}

func GetUTXO(c *gin.Context) {

	address := c.Param("address")

	utxoArray, err  := btcUtils.GetUTXO(address)
	if err != nil{
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	response := new(responses.UTXOResponse)

	response.Utxo = utxoArray

	c.JSON(http.StatusOK, response)
}

//func GetBalances(c *gin.Context) {
//
//	type Request struct {
//		AddressesArray []string `json:"addressesArray"`
//	}
//
//	request := new(Request)
//
//	err := c.BindJSON(&request)
//	if err != nil {
//		log.Println(err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
//		return
//	}
//
//	var wg sync.WaitGroup
//
//	balances := multiBalance.New()
//
//	for i := 0; i < len(request.AddressesArray); i++ {
//		wg.Add(1)
//		go multiBalance.BtcWorker(&wg, request.AddressesArray[i], balances)
//	}
//	wg.Wait()
//
//	response := new(responses.BalancesResponse)
//	response.Balances = balances.Result
//
//	c.JSON(http.StatusOK, response)
//}
