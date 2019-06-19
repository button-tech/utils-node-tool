package handlers

import (
	"github.com/button-tech/utils-node-tool/shared/btcUtils"
	"github.com/button-tech/utils-node-tool/shared/responseModels"
	"github.com/qiangxue/fasthttp-routing"
)

func GetBalance(c *routing.Context) error {

	address := c.Param("address")

	response := new(responses.BalanceResponse)

	balance, err := btcUtils.GetBtcBlockChainBalance(address)
	if err != nil {
		return err
	}

	response.Balance = balance

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	return nil
}

func GetUTXO(c *routing.Context) error {

	address := c.Param("address")

	utxoArray, err := btcUtils.GetUTXO(address)
	if err != nil {
		return err
	}

	response := new(responses.UTXOResponse)

	response.Utxo = utxoArray

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	return nil
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
