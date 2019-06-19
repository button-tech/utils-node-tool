package handlers

import (
	"math"
	"net/http"
	"sync"

	"log"

	"context"
	"math/big"
	"os"

	"github.com/button-tech/utils-node-tool/shared/abi"
	"github.com/button-tech/utils-node-tool/shared/db"
	"github.com/button-tech/utils-node-tool/shared/multiBalance"
	"github.com/button-tech/utils-node-tool/shared/responseModels"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/onrik/ethrpc"
	"golang.org/x/crypto/sha3"
)

func GetBalance(c *gin.Context) {

	var ethClient = ethrpc.New(os.Getenv("eth-api"))

	balance, err := ethClient.EthGetBalance(c.Param("address"), "latest")
	if err != nil {
		reserveNode, err := db.GetEndpoint("eth")
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		ethClient = ethrpc.New(reserveNode)

		result, err := ethClient.EthGetBalance(c.Param("address"), "latest")
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

func GetTxFee(c *gin.Context) {

	ethClient := ethrpc.New(os.Getenv("eth-api"))

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

func GetGasPrice(c *gin.Context) {

	ethClient := ethrpc.New(os.Getenv("eth-api"))

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

func GetTokenBalance(c *gin.Context) {

	address := c.Param("address")

	smartContractAddress := c.Param("sc-address")

	ethClient, err := ethclient.Dial(os.Getenv("eth-api"))
	if err != nil {
		endPoint, err := db.GetEndpoint("eth")
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		ethClient, err = ethclient.Dial(endPoint)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}

	instance, err := abi.NewToken(common.HexToAddress(smartContractAddress), ethClient)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	balance, err := instance.BalanceOf(nil, common.HexToAddress(address))
	if err != nil {
		endPoint, err := db.GetEndpoint("eth")
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		ethClient, err = ethclient.Dial(endPoint)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		instance, err = abi.NewToken(common.HexToAddress(smartContractAddress), ethClient)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		balance, err = instance.BalanceOf(nil, common.HexToAddress(address))
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}

	response := new(responses.BalanceResponse)

	response.Balance = balance.String()

	c.JSON(http.StatusOK, response)
}

func GetEstimateGas(c *gin.Context) {

	txData := struct {
		ToAddress    string `json:"toAddress"`
		TokenAddress string `json:"tokenAddress"`
		Amount       string `json:"amount"`
	}{}

	err := c.BindJSON(&txData)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	toAddress := common.HexToAddress(txData.ToAddress)
	tokenAddress := common.HexToAddress(txData.TokenAddress)

	amount := new(big.Int)
	amount.SetString(txData.Amount, 10)

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	_, err = hash.Write(transferFnSignature)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	methodID := hash.Sum(nil)[:4]

	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	ethClient, err := ethclient.Dial(os.Getenv("eth-api"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	gasLimit, err := ethClient.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	})

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"gasLimit": gasLimit})
}

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
		go multiBalance.EthWorker(&wg, req.AddressesArray[i], balances)
	}
	wg.Wait()

	response := new(responses.BalancesResponse)
	response.Balances = balances.Result

	c.JSON(http.StatusOK, response)
}

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
