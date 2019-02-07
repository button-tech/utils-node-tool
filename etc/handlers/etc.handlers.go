package handlers

import (
	"context"
	"github.com/button-tech/utils-node-tool/etc/handlers/multiBalance"
	"github.com/button-tech/utils-node-tool/etc/handlers/responseModels"
	. "github.com/button-tech/utils-node-tool/etc/handlers/storage"
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"net/http"
	"sync"
)

var ctx = context.Background()

// @Summary ETC balance of account
// @Description return balance of account in ETC for specific node
// @Produce  application/json
// @Param   address        path    string     true        "address"
// @Success 200 {array} responses.BalanceResponse
// @Router /etc/balance/{address} [get]
// GetBalance return balance of account in ETC for specific node
func GetBalance(c *gin.Context) {

	balance, err := EtcClient.EthGetBalance(c.Param("address"), "latest")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	response := new(responses.BalanceResponse)

	response.Balance = balance.String()

	c.JSON(http.StatusOK, response)
}

// @Summary return Amount of ETC that you need to send a transaction
// @Description return Amount of ETC that you need to send a transaction
// @Produce  application/json
// @Success 200 {array} responses.TransactionFeeResponse
// @Router /etc/transactionFee [get]
// GetTxFee return Amount of ETC that you need to send a transaction
func GetTxFee(c *gin.Context) {

	gasPrice, err := EtcClient.EthGasPrice()

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	fee := float64(gasPrice.Int64()*21000) / math.Pow(10, 18)

	response := new(responses.TransactionFeeResponse)
	response.Fee = fee

	c.JSON(http.StatusOK, response)
}

// @Summary return gas price of specific node
// @Description return Amount of ETC that you need to send a transaction
// @Produce  application/json
// @Success 200 {array} responses.GasPriceResponse
// @Router /etc/gasPrice [get]
// GetGasPrice return gas price of specific node
func GetGasPrice(c *gin.Context) {

	gasPrice, err := EtcClient.EthGasPrice()

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	response := new(responses.GasPriceResponse)
	response.GasPrice = gasPrice.Int64()

	c.JSON(http.StatusOK, response)
}

// @Summary ETC balance of accounts by list
// @Description return balances of accounts in ETC
// @Produce  application/json
// @Param addressesArray     body string true "addressesArray"
// @Success 200 {array} responses.BalancesResponse
// @Router /etc/balances [post]
// GetBalanceForMultipleAdresses return balances of accounts in ETC
func GetBalances(c *gin.Context) {

	type Request struct {
		AddressesArray []string `json:"addressesArray"`
	}

	req := new(Request)

	var balances multiBalance.Balances

	c.BindJSON(&req)

	var wg sync.WaitGroup

	for i := 0; i < len(req.AddressesArray); i++ {
		wg.Add(1)
		go multiBalance.Worker(&wg, req.AddressesArray[i], &balances)
	}
	wg.Wait()

	c.JSON(http.StatusOK, balances.Result)
}

//not Working yet on production
// func SendTX(c *gin.Context) {

// 	etcClient, err := ethclient.Dial(etcURL)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	type DataToSend struct {
// 		RawTx string `json:"RawTx"`
// 	}
// 	var data DataToSend

// 	c.BindJSON(&data)

// 	fmt.Println("RawTx: " + data.RawTx)

// 	rawTxBytes, err := hex.DecodeString(data.RawTx)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	tx := new(types.Transaction)
// 	rlp.DecodeBytes(rawTxBytes, &tx)

// 	err = etcClient.SendTransaction(context.Background(), tx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("tx sent: %s", tx.Hash().Hex())

// 	c.JSON(http.StatusOK, gin.H{"response": tx.Hash().Hex()})
// }
