package handlers

import (
	"context"
	"math"
	"net/http"
	"sync"

	"github.com/button-tech/utils-node-tool/eth/abi"
	"github.com/button-tech/utils-node-tool/eth/handlers/multiBalance"
	"github.com/button-tech/utils-node-tool/eth/handlers/responseModels"
	"github.com/button-tech/utils-node-tool/eth/handlers/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/onrik/ethrpc"
	"log"
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

	var EthClient = ethrpc.New(storage.EthNodeAddress.Address)

	balance, err := EthClient.EthGetBalance(c.Param("address"), "latest")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	response := new(responses.BalanceResponse)

	response.Balance = balance.String()

	c.JSON(http.StatusOK, response)
}

// @Summary return Amount of ETH that you need to send a transaction
// @Description return Amount of ETH that you need to send a transaction
// @Produce  application/json
// @Success 200 {array} responses.TransactionFeeResponse
// @Router /eth/transactionFee [get]
// GetTxFee return Amount of ETH that you need to send a transaction
func GetTxFee(c *gin.Context) {

	gasPrice, err := storage.EthClient.EthGasPrice()

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
// @Description return Amount of ETH that you need to send a transaction
// @Produce  application/json
// @Success 200 {array} responses.GasPriceResponse
// @Router /eth/gasPrice [get]
// GetGasPrice return gas price of specific node
func GetGasPrice(c *gin.Context) {

	gasPrice, err := storage.EthClient.EthGasPrice()

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	response := new(responses.GasPriceResponse)
	response.GasPrice = gasPrice.Int64()

	c.JSON(http.StatusOK, response)
}

// @Summary return balance of specific token in ETH node
// @Description return balance of specific token in ETH node
// @Produce  application/json
// @Param   address        path    string     true        "address"
// @Param   sc-address        path    string     true        "sc-address"
// @Success 200 {array} responses.BalanceResponse
// @Router /eth/tokenBalance/{sc-address}/{address} [get]
// GetTokenBalance return Amount of ETH ERC20 token
func GetTokenBalance(c *gin.Context) {

	address := c.Param("address")

	smartContractAddress := c.Param("sc-address")

	ethClient, err := ethclient.Dial(storage.EthNodeAddress.Address)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	instance, err := abi.NewToken(common.HexToAddress(smartContractAddress), ethClient)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	balance, err := instance.BalanceOf(nil, common.HexToAddress(address))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 500})
		return
	}

	response := new(responses.BalanceResponse)

	response.Balance = balance.String()

	c.JSON(http.StatusOK, response)
}

// @Summary ETH balance of accounts by list
// @Description return balances of accounts in ETH
// @Produce  application/json
// @Param addressesArray     body string true "addressesArray"
// @Success 200 {array} responses.BalancesResponse
// @Router /eth/balances [post]
// GetBalanceForMultipleAdresses return balances of accounts in ETH
func GetBalances(c *gin.Context) {

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

// @Summary ETH ERC-20 tokens balance of account by list of smart contracts
// @Description return tokens balances of account
// @Produce  application/json
// @Param addressesArray     body string true "addressesArray"
// @Success 200 {array} responses.BalancesResponse
// @Router /eth/tokenBalances [post]
// GetBalanceForMultipleAdresses return tokens balances of account
func GetTokenBalances(c *gin.Context) {

	type Request struct {
		OwnerAddress   string   `json:"ownerAddress""`
		SmartAddresses []string `json:"smartAddresses"`
	}

	req := new(Request)

	balances := multiBalance.New()

	c.BindJSON(&req)

	var wg sync.WaitGroup

	for i := 0; i < len(req.SmartAddresses); i++ {
		wg.Add(1)
		go multiBalance.TokenWorker(&wg, req.OwnerAddress, req.SmartAddresses[i], balances)
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
