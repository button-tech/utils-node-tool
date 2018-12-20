package handlers

import (
	"context"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/button-tech/utils-node-tool/etc/handlers/responseModels"
	"github.com/gin-gonic/gin"
	"github.com/onrik/ethrpc"
)

var (
	etcURL = os.Getenv("ETC_NODE")
)

var (
	ctx       = context.Background()
	etcClient = ethrpc.New(etcURL)
)

// @Summary ETC balance of account
// @Description return balance of account in ETC for specific node
// @Produce  application/json
// @Param   address        path    string     true        "address"
// @Success 200 {array} responses.BalanceResponse
// @Router /etc/balance/{address} [get]
// GetBalance return balance of account in ETC for specific node
func GetBalance(c *gin.Context) {

	balance, err := etcClient.EthGetBalance(c.Param("address"), "latest")
	if err != nil {
		log.Println(err)
	}
	floatBalance, _ := strconv.ParseFloat(balance.String(), 64)

	ethBalance := floatBalance / math.Pow(10, 18)

	response := new(responses.BalanceResponse)
	response.Balance = ethBalance

	c.JSON(http.StatusOK, response)
}

// @Summary return Amount of ETC that you need to send a transaction
// @Description return Amount of ETC that you need to send a transaction
// @Produce  application/json
// @Success 200 {array} responses.TransactionFeeResponse
// @Router /etc/transactionFee [get]
// GetBalance return Amount of ETC that you need to send a transaction
func GetTxFee(c *gin.Context) {

	gasPrice, err := etcClient.EthGasPrice()

	if err != nil {
		log.Println(err)
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
// GetBalance return gas price of specific node
func GetGasPrice(c *gin.Context) {

	gasPrice, err := etcClient.EthGasPrice()

	if err != nil {
		log.Println(err)
	}

	response := new(responses.GasPriceResponse)
	response.GasPrice = gasPrice.Int64()

	c.JSON(http.StatusOK, response)
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
