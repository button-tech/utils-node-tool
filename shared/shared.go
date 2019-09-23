package shared

import (
	"context"
	"errors"
	"github.com/button-tech/utils-node-tool/shared/db"
	"github.com/button-tech/utils-node-tool/shared/requests"
	"github.com/button-tech/utils-node-tool/shared/responses"
	"github.com/button-tech/utils-node-tool/utils-for-endpoints/storage"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imroc/req"
	"log"
	"os"
	"strconv"
	"strings"
	"encoding/hex"
)

func GetEstimateGas(req *requests.EthEstimateGasRequest) (uint64, error) {

	ethClient, err := ethclient.Dial(storage.EndpointForReq.Get())
	if err != nil {
		return 0, err
	}

	address := common.HexToAddress(req.ContractAddress)

	data, err := hex.DecodeString(req.Data)
	if err != nil{
		return 0, err
	}

	gasLimit, err := ethClient.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &address,
		Data: data,
	})

	if err != nil {
		return 0, err
	}

	return gasLimit, nil
}

func GetUtxo(address string) ([]responses.UTXO, error) {

	currency := os.Getenv("BLOCKCHAIN")

	var requestUrl, endPoint string

	var err error

	if currency != "bch" {
		endPoint, err = storage.GetEndpoint(currency)
		if err != nil {
			return nil, err
		}
	}

	switch currency {
	case "btc":
		requestUrl = endPoint + "/addr/" + address + "/utxo"
	case "bch":
		requestUrl = "https://rest.bitbox.earth/v1/address/utxo/" + address
	case "ltc":
		requestUrl = endPoint + "/api/addr/" + address + "/utxo"
	}

	utxos, err := req.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	if utxos.Response().StatusCode != 200 {
		return nil, errors.New("Bad request")
	}

	var utxoArray []responses.UTXO

	err = utxos.ToJSON(&utxoArray)
	if err != nil {
		return nil, err
	}

	return utxoArray, nil
}

func GetEthBasedBlockNumber(currency, addr string) (int64, error) {
	header := req.Header{
		"Content-Type": "application/json",
	}

	params := strings.NewReader("{\n\"jsonrpc\":\"2.0\",\n\"method\":\"eth_getBlockByNumber\",\n\"params\":[\"latest\", false],\n\"id\":1\n}")

	resp, err := req.Post(addr, header, params)

	if err != nil || resp.Response().StatusCode != 200 {
		err := DeleteEntry(currency, addr)
		if err != nil {
			return 0, err
		}
		return 0, nil
	}

	info := struct {
		Result struct {
			Number string `json:"number"`
		}
	}{}

	err = resp.ToJSON(&info)
	if err != nil {
		return 0, err
	}

	if len(info.Result.Number) == 0 {
		return 0, errors.New("Bad request")
	}

	hexNumber := []byte(info.Result.Number)

	intNumber, err := strconv.ParseInt(string(hexNumber[2:]), 16, 64)
	if err != nil {
		return 0, err
	}

	return intNumber, nil
}

func GetUtxoBasedBlockNumber(currency, addr string) (int64, error) {
	var url string

	switch currency {
	case "btc":
		url = "/sync"
	case "ltc":
		url = "/api/sync"
	}

	res, err := req.Get(addr + url)
	if err != nil || res.Response().StatusCode != 200 {
		err := DeleteEntry(currency, addr)
		if err != nil {
			return 0, err
		}
		log.Println("Status code:" + strconv.Itoa(res.Response().StatusCode))
	}

	info := struct {
		BlockChainHeight int64 `json:"blockChainHeight"`
	}{}

	err = res.ToJSON(&info)
	if err != nil {
		return 0, err
	}

	return info.BlockChainHeight, nil
}

func Max(array []int64) int64 {
	var max int64 = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
	}
	return max
}

func DeleteEntry(currency, address string) error {
	err := db.AddToStoppedList(currency, address)
	if err != nil {
		return err
	}

	log.Printf("Add to stopped list %s: %s", currency, address)

	return nil
}
