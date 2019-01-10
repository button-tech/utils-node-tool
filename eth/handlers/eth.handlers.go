package handlers

import (
	"context"
	"log"
	"math"
	"net/http"
	"strconv"
	"sync"

	"github.com/button-tech/utils-node-tool/eth/abi"
	"github.com/button-tech/utils-node-tool/eth/handlers/responseModels"
	"github.com/button-tech/utils-node-tool/eth/multiBalance"
	"github.com/button-tech/utils-node-tool/eth/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
)

var ctx = context.Background()

// @Summary ETH balance of account
// @Description return balance of account in ETH for specific node
// @Produce  application/json
// @Param   address        path    string     true        "address"
// @Success 200 {array} responses.BalanceResponse
// @Router /eth/balance/{address} [get]
// GetBalance return balance of account in ETH for specific node
func GetBalance(c *gin.Context) {

	balance, err := storage.EthClient.EthGetBalance(c.Param("address"), "latest")
	if err != nil {
		log.Println(err)
	}
	floatBalance, _ := strconv.ParseFloat(balance.String(), 64)

	ethBalance := floatBalance / math.Pow(10, 18)

	response := new(responses.BalanceResponse)
	response.Balance = ethBalance

	c.JSON(http.StatusOK, response)
}

// @Summary return Amount of ETH that you need to send a transaction
// @Description return Amount of ETH that you need to send a transaction
// @Produce  application/json
// @Success 200 {array} responses.TransactionFeeResponse
// @Router /eth/transactionFee [get]
// GetBalance return Amount of ETH that you need to send a transaction
func GetTxFee(c *gin.Context) {

	gasPrice, err := storage.EthClient.EthGasPrice()

	if err != nil {
		log.Println(err)
	}

	fee := float64(gasPrice.Int64()*21000) / math.Pow(10, 18)

	response := new(responses.TransactionFeeResponse)
	response.Fee = fee

	c.JSON(http.StatusOK, response)
}

// @Summary return gas price of specific node
// @Description return Amount of ETH that you need to send a transaction
// @Produce  application/json
// @Success 200 {array} responses.GasPriceResponse
// @Router /eth/gasPrice [get]
// GetBalance return gas price of specific node
func GetGasPrice(c *gin.Context) {

	gasPrice, err := storage.EthClient.EthGasPrice()

	if err != nil {
		log.Println(err)
	}

	response := new(responses.GasPriceResponse)
	response.GasPrice = gasPrice.Int64()

	c.JSON(http.StatusOK, response)
}

// @Summary return balance of specific token in ETH node
// @Description return balance of specific token in ETH node
// @Produce  application/json
// @Param   address        path    string     true        "address"
// @Param   token        path    string     true        "token"
// @Success 200 {array} responses.TokenBalanceResponse
// @Router /eth/tokenBalance/{token}/{address} [get]
// GetBalance return Amount of ETH ERC20 token
func GetTokenBalance(c *gin.Context) {

	address := c.Param("address")

	token := c.Param("token")

	ethClient, err := ethclient.Dial(storage.EthURL)
	if err != nil {
		log.Println(err)
	}

	instance, err := abi.NewToken(common.HexToAddress(storage.TokensAddresses[token]), ethClient)
	if err != nil {
		log.Println(err)
	}

	localBalance, err := instance.BalanceOf(nil, common.HexToAddress(address))
	if err != nil {
		log.Println(err)
	}

	floatTokenBalance, _ := strconv.ParseFloat(localBalance.String(), 64)

	tokenBalance := floatTokenBalance / math.Pow(10, 18)

	response := new(responses.TokenBalanceResponse)
	response.TokenBalance = tokenBalance

	c.JSON(http.StatusOK, response)
}

// @Summary ETH balance of accounts by list
// @Description return balances of accounts in ETH
// @Produce  application/json
// @Param addressesArray     body string true "addressesArray"
// @Success 200 {array} responses.BalancesResponse
// @Router /eth/balances [post]
// GetBalanceForMultipleAdresses return balances of accounts in ETH
func GetBalanceForMultipleAdresses(c *gin.Context) {

	type Request struct {
		AddressesArray []string `json:"addressesArray"`
	}

	req := new(Request)

	balances := multiBalance.New()

	c.BindJSON(&req)

	var wg sync.WaitGroup

	for i := 0; i < len(req.AddressesArray); i++ {
		wg.Add(1)
		go multiBalance.Worker(&wg, req.AddressesArray[i], balances)
	}
	wg.Wait()

	response := new(responses.BalancesResponse)
	response.Balances = balances.Result

	c.JSON(http.StatusOK, response)
}

//not Working yet on production
// func SendTX(c *gin.Context) {

// 	ethClient, err := ethclient.Dial(ethURL)
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

// 	err = ethClient.SendTransaction(context.Background(), tx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("tx sent: %s", tx.Hash().Hex())

// 	c.JSON(http.StatusOK, gin.H{"response": tx.Hash().Hex()})
// }
