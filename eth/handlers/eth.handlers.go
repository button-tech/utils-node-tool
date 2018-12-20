package handlers

import (
	"context"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/button-tech/utils-node-tool/eth/abi"
	"github.com/button-tech/utils-node-tool/eth/handlers/responseModels"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/onrik/ethrpc"
)

var (
	ethURL = os.Getenv("ETH_NODE")
)

var (
	ctx = context.Background()

	ethClient = ethrpc.New(ethURL)

	tokensAddresses = map[string]string{

		"bix":  "0xb3104b4b9da82025e8b9f8fb28b3553ce2f67069",
		"btm":  "0xcb97e65f07da24d46bcdd078ebebd7c6e6e3d750",
		"omg":  "0xd26114cd6EE289AccF82350c8d8487fedB8A0C07",
		"elf":  "0xbf2179859fc6d5bee9bf9158632dc51678a4100e",
		"bnb":  "0xB8c77482e45F1F44dE1745F52C74426C631bDD52",
		"tusd": "0x8dd5fbce2f6a956c3022ba3663759011dd51e73e",
		"knc":  "0xdd974d5c2e2928dea5f71b9825b8b646686bd200",
		"zrx":  "0xe41d2489571d322189246dafa5ebde1f4699f498",
		"rep":  "0x1985365e9f78359a9B6AD760e32412f4a445E862",
		"gnt":  "0xa74476443119A942dE498590Fe1f2454d7D4aC0d",
	}
)

// @Summary ETH balance of account
// @Description return balance of account in ETH for specific node
// @Produce  application/json
// @Param   address        path    string     true        "address"
// @Success 200 {array} responses.BalanceResponse
// @Router /eth/balance/{address} [get]
// GetBalance return balance of account in ETH for specific node
func GetBalance(c *gin.Context) {

	balance, err := ethClient.EthGetBalance(c.Param("address"), "latest")
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

	gasPrice, err := ethClient.EthGasPrice()

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

	gasPrice, err := ethClient.EthGasPrice()

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

	ethClient, err := ethclient.Dial(ethURL)
	if err != nil {
		log.Println(err)
	}

	instance, err := abi.NewToken(common.HexToAddress(tokensAddresses[token]), ethClient)
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
