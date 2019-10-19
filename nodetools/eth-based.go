package nodetools

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/button-tech/utils-node-tool/nodetools/storage"
	"github.com/button-tech/utils-node-tool/requests"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imroc/req"
	"strconv"
	"strings"
)

func GetEstimateGas(req *requests.EthEstimateGasRequest) (uint64, error) {

	ethClient, err := ethclient.Dial(storage.EndpointForReq.Get())
	if err != nil {
		return 0, err
	}

	address := common.HexToAddress(req.ContractAddress)

	data, err := hex.DecodeString(req.Data)
	if err != nil {
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

func GetEthBasedBlockNumber(addr string) (int64, error) {
	header := req.Header{
		"Content-Type": "application/json",
	}

	params := strings.NewReader("{\n\"jsonrpc\":\"2.0\",\n\"method\":\"eth_getBlockByNumber\",\n\"params\":[\"latest\", false],\n\"id\":1\n}")

	resp, err := req.Post(addr, header, params)
	if err != nil {
		return 0, err
	}

	if resp.Response().StatusCode != 200 {
		return 0, errors.New("Bad request")
	}

	var info requests.EthBasedBlocksHeight

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
