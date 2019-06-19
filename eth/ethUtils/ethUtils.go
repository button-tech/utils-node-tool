package ethUtils

import (
	"github.com/onrik/ethrpc"
	"github.com/button-tech/utils-node-tool/shared/db"
	"os"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/button-tech/utils-node-tool/shared/abi"
	"math/big"
	"golang.org/x/crypto/sha3"
	"github.com/ethereum/go-ethereum"
	"context"
)

type TxData struct {
	ToAddress    string `json:"toAddress"`
	TokenAddress string `json:"tokenAddress"`
	Amount       string `json:"amount"`
}

func GetBalance(address string) (string, error){

	var ethClient = ethrpc.New(os.Getenv("main-api"))

	balance, err := ethClient.EthGetBalance(address, "latest")
	if err != nil {
		reserveNode, err := db.GetEndpoint("blockChain")
		if err != nil {
			return "", err
		}

		ethClient = ethrpc.New(reserveNode)

		result, err := ethClient.EthGetBalance(address, "latest")
		if err != nil {
			return "", err
		}

		balance = result
	}

	return balance.String(), nil
}


func GetTokenBalance(userAddress, smartContractAddress string)(string, error){
	ethClient, err := ethclient.Dial(os.Getenv("main-api"))
	if err != nil {
		endPoint, err := db.GetEndpoint("blockChain")
		if err != nil {
			return "", err
		}

		ethClient, err = ethclient.Dial(endPoint)
		if err != nil {
			return "", err
		}
	}

	instance, err := abi.NewToken(common.HexToAddress(smartContractAddress), ethClient)
	if err != nil {
		return "", err
	}

	balance, err := instance.BalanceOf(nil, common.HexToAddress(userAddress))
	if err != nil {
		endPoint, err := db.GetEndpoint("main-api")
		if err != nil {
			return "", err
		}

		ethClient, err = ethclient.Dial(endPoint)
		if err != nil {
			return "", err
		}

		instance, err = abi.NewToken(common.HexToAddress(smartContractAddress), ethClient)
		if err != nil {
			return "", err
		}

		balance, err = instance.BalanceOf(nil, common.HexToAddress(userAddress))
		if err != nil {
			return "", err
		}
	}

	return balance.String(), nil
}

func GetEstimateGas(txData *TxData)(uint64, error){

	toAddress := common.HexToAddress(txData.ToAddress)

	tokenAddress := common.HexToAddress(txData.TokenAddress)

	amount := new(big.Int)
	amount.SetString(txData.Amount, 10)

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()

	_, err := hash.Write(transferFnSignature)
	if err != nil {
		return 0, err
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
		return 0, err
	}

	gasLimit, err := ethClient.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	})

	if err != nil {
		return 0, err
	}

	return gasLimit, nil
}