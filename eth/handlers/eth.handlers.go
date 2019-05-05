package handlers

import (
	"math"
	"net/http"
	"sync"

	"log"

	"github.com/button-tech/utils-node-tool/db"
	"github.com/button-tech/utils-node-tool/eth/abi"
	"github.com/button-tech/utils-node-tool/eth/handlers/multiBalance"
	"github.com/button-tech/utils-node-tool/eth/handlers/responseModels"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/onrik/ethrpc"
)

// @Summary ETH balance of account
// @Description return balance of account in ETH for specific node
// @Produce  application/json
// @Param   address        path    string     true        "address"
// @Success 200 {array} responses.BalanceResponse
// @Router /eth/balance/{address} [get]
// GetBalance return balance of account in ETH for specific node
func GetBalance(c *gin.Context) {
	endPoint, err := db.GetEndpoint("eth")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var EthClient = ethrpc.New(endPoint)

	balance, err := EthClient.EthGetBalance(c.Param("address"), "latest")

	if err != nil {

		reserveNode, err := db.GetReserveHost("eth")
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		EthClient = ethrpc.New(reserveNode)

		result, err := EthClient.EthGetBalance(c.Param("address"), "latest")
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		balance = result
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

	endPoint, err := db.GetEndpoint("eth")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ethClient := ethrpc.New(endPoint)

	gasPrice, err := ethClient.EthGasPrice()

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
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

	endPoint, err := db.GetEndpoint("eth")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ethClient := ethrpc.New(endPoint)

	gasPrice, err := ethClient.EthGasPrice()

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
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

	endPoint, err := db.GetEndpoint("eth")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ethClient, err := ethclient.Dial(endPoint)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	instance, err := abi.NewToken(common.HexToAddress(smartContractAddress), ethClient)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	balance, err := instance.BalanceOf(nil, common.HexToAddress(address))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
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

	err := c.BindJSON(&req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

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
		OwnerAddress   string   `json:"ownerAddress"`
		SmartAddresses []string `json:"smartAddresses"`
	}

	req := new(Request)

	balances := multiBalance.New()

	err := c.BindJSON(&req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

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
