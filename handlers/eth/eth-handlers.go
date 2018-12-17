package eth

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"./abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gin-gonic/gin"
	"github.com/onrik/ethrpc"
)



var (
	ethURL =  os.Getenv("ETH_NODE")
	etcURL =  os.Getenv("ETC_NODE")
)

var (
	ctx = context.Background()

	ethClient = ethrpc.New(ethURL)

	etcClient = ethrpc.New(etcURL)

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

func GetBalance(c *gin.Context) {

	client := ClientRpc(c.Param("eth"))

	balance, err := client.EthGetBalance(c.Param("address"), "latest")
	if err != nil {
		log.Println(err)
	}
	floatBalance, _ := strconv.ParseFloat(balance.String(), 64)

	ethBalance := floatBalance / math.Pow(10, 18)

	c.JSON(http.StatusOK, gin.H{"balance": ethBalance})

	return
}

func GetTxFee(c *gin.Context) {

	client := ClientRpc(c.Param("eth"))

	gasPrice, err := client.EthGasPrice()

	if err != nil {
		log.Println(err)
	}

	fee := float64(gasPrice.Int64()*21000) / math.Pow(10, 18)

	c.JSON(http.StatusOK, gin.H{"fee": fee})

	return
}

func GetGasPrice(c *gin.Context) {

	client := ClientRpc(c.Param("eth"))

	gasPrice, err := client.EthGasPrice()

	if err != nil {
		log.Println(err)
	}
	
	c.JSON(http.StatusOK, gin.H{"gas_price": gasPrice.Int64()})
}

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

	balance, err := instance.BalanceOf(nil, common.HexToAddress(address))

	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
	return
}

// Send Tx
func SendTX(c *gin.Context) {

	client := ClientDial(c.Param("eth"))

	type DataToSend struct {
		Raw_tx string `json:"raw_tx"`
	}
	var data DataToSend

	c.BindJSON(&data)

	fmt.Println("RawTx: " + data.Raw_tx)

	rawTxBytes, err := hex.DecodeString(data.Raw_tx)
	if err != nil {
		fmt.Println(err)
	}

	tx := new(types.Transaction)
	rlp.DecodeBytes(rawTxBytes, &tx)

	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tx sent: %s", tx.Hash().Hex())

	c.JSON(http.StatusOK, gin.H{"response": tx.Hash().Hex()})
}


// get client
func ClientRpc(param string) *ethrpc.EthRPC {
	if param == "etc" {
		return etcClient
	}
	return ethClient
}

func ClientDial(param string) *ethclient.Client{
	if param == "etc"{
		client, err := ethclient.Dial(etcURL)
		if err != nil {
			log.Println(err)
		}
		return client
	}

	client, err := ethclient.Dial(ethURL)
	if err!= nil{
		log.Println(err)
	}
	return client
}