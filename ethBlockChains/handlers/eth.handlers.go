package handlers

import (
	"log"
	"math"

	"context"
	"math/big"
	"os"

	"encoding/json"
	"github.com/button-tech/utils-node-tool/shared/abi"
	"github.com/button-tech/utils-node-tool/shared/db"
	"github.com/button-tech/utils-node-tool/shared/responseModels"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/onrik/ethrpc"
	"github.com/qiangxue/fasthttp-routing"
	"golang.org/x/crypto/sha3"
)

func GetBalance(c *routing.Context) error {

	var ethClient = ethrpc.New(os.Getenv("main-api"))

	balance, err := ethClient.EthGetBalance(c.Param("address"), "latest")
	if err != nil {
		reserveNode, err := db.GetEndpoint("blockChain")
		if err != nil {
			log.Println(err)
			return err
		}

		ethClient = ethrpc.New(reserveNode)

		result, err := ethClient.EthGetBalance(c.Param("address"), "latest")
		if err != nil {
			log.Println(err)
			return err
		}

		balance = result
	}

	response := new(responses.BalanceResponse)

	response.Balance = balance.String()

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	return nil
}

func GetTxFee(c *routing.Context) error {

	ethClient := ethrpc.New(os.Getenv("main-api"))

	gasPrice, err := ethClient.EthGasPrice()

	if err != nil {
		log.Println(err)
		return err
	}

	fee := float64(gasPrice.Int64()*21000) / math.Pow(10, 18)

	response := new(responses.TransactionFeeResponse)

	response.Fee = fee

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	return nil
}

func GetGasPrice(c *routing.Context) error {

	ethClient := ethrpc.New(os.Getenv("main-api"))

	gasPrice, err := ethClient.EthGasPrice()

	if err != nil {
		log.Println(err)
		return err
	}

	response := new(responses.GasPriceResponse)

	response.GasPrice = gasPrice.Int64()

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	return nil
}

func GetTokenBalance(c *routing.Context) error {

	address := c.Param("address")

	smartContractAddress := c.Param("sc-address")

	ethClient, err := ethclient.Dial(os.Getenv("main-api"))
	if err != nil {
		endPoint, err := db.GetEndpoint("blockChain")
		if err != nil {
			log.Println(err)
			return err
		}

		ethClient, err = ethclient.Dial(endPoint)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	instance, err := abi.NewToken(common.HexToAddress(smartContractAddress), ethClient)
	if err != nil {
		log.Println(err)
		return err
	}

	balance, err := instance.BalanceOf(nil, common.HexToAddress(address))
	if err != nil {
		endPoint, err := db.GetEndpoint("main-api")
		if err != nil {
			log.Println(err)
			return err
		}

		ethClient, err = ethclient.Dial(endPoint)
		if err != nil {
			log.Println(err)
			return err
		}

		instance, err = abi.NewToken(common.HexToAddress(smartContractAddress), ethClient)
		if err != nil {
			log.Println(err)
			return err
		}

		balance, err = instance.BalanceOf(nil, common.HexToAddress(address))
		if err != nil {
			log.Println(err)
			return err
		}
	}

	response := new(responses.BalanceResponse)

	response.Balance = balance.String()

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	return nil
}

func GetEstimateGas(c *routing.Context) error {

	txData := struct {
		ToAddress    string `json:"toAddress"`
		TokenAddress string `json:"tokenAddress"`
		Amount       string `json:"amount"`
	}{}

	if err := json.Unmarshal(c.PostBody(), &txData); err != nil {
		log.Println(err)
		return err
	}

	toAddress := common.HexToAddress(txData.ToAddress)
	tokenAddress := common.HexToAddress(txData.TokenAddress)

	amount := new(big.Int)
	amount.SetString(txData.Amount, 10)

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	_, err := hash.Write(transferFnSignature)
	if err != nil {
		log.Println(err)
		return err
	}
	methodID := hash.Sum(nil)[:4]

	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	ethClient, err := ethclient.Dial(os.Getenv("main-api"))
	if err != nil {
		log.Println(err)
		return err
	}

	gasLimit, err := ethClient.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	})

	if err != nil {
		log.Println(err)
		return err
	}

	response := new(responses.GasLimitResponse)
	response.GasLimit = gasLimit

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
//	req := new(Request)
//
//	balances := multiBalance.New()
//
//	err := c.BindJSON(&req)
//	if err != nil {
//		log.Println(err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
//		return
//	}
//
//	var wg sync.WaitGroup
//
//	for i := 0; i < len(req.AddressesArray); i++ {
//		wg.Add(1)
//		go multiBalance.EthWorker(&wg, req.AddressesArray[i], balances)
//	}
//	wg.Wait()
//
//	response := new(responses.BalancesResponse)
//	response.Balances = balances.Result
//
//	c.JSON(http.StatusOK, response)
//}
//
//func GetTokenBalances(c *gin.Context) {
//
//	type Request struct {
//		OwnerAddress   string   `json:"ownerAddress"`
//		SmartAddresses []string `json:"smartAddresses"`
//	}
//
//	req := new(Request)
//
//	balances := multiBalance.New()
//
//	err := c.BindJSON(&req)
//	if err != nil {
//		log.Println(err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
//		return
//	}
//
//	var wg sync.WaitGroup
//
//	for i := 0; i < len(req.SmartAddresses); i++ {
//		wg.Add(1)
//		go multiBalance.TokenWorker(&wg, req.OwnerAddress, req.SmartAddresses[i], balances)
//	}
//	wg.Wait()
//
//	response := new(responses.BalancesResponse)
//	response.Balances = balances.Result
//
//	c.JSON(http.StatusOK, response)
//}
