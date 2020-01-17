package handlers

import (
	"encoding/json"
	"errors"
	"github.com/button-tech/logger"
	"github.com/button-tech/utils-node-tool/nodetools"
	b "github.com/button-tech/utils-node-tool/nodetools"
	"github.com/button-tech/utils-node-tool/nodetools/storage"
	"github.com/button-tech/utils-node-tool/types/requests"
	"github.com/button-tech/utils-node-tool/types/responses"
	"github.com/imroc/req"
	"github.com/qiangxue/fasthttp-routing"
	"os"
	"time"
)

func GetBalance(c *routing.Context) error {

	start := time.Now()

	address := c.Param("address")

	balance, err := b.GetUtxoBasedBalance(address)
	if err != nil {
		logger.HandlerError("GetBalance", err)
		return err
	}

	response := responses.BalanceResponse{Balance: balance}

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	logger.LogRequest(time.Since(start), os.Getenv("BLOCKCHAIN"), "GetBalance", false)

	return nil
}

func GetUtxo(c *routing.Context) error {

	start := time.Now()

	address := c.Param("address")

	utxoArray, err := nodetools.GetUtxo(address)
	if err != nil {
		logger.HandlerError("GetUtxo", err)
		return err
	}

	response := responses.UTXOResponse{Utxo: utxoArray}

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	logger.LogRequest(time.Since(start), os.Getenv("BLOCKCHAIN"), "GetUtxo", false)

	return nil
}

func SendRawTx(c *routing.Context) error {
	start := time.Now()

	var (
		rawTx  requests.RawTransaction
		result responses.TransactionResult
	)

	if err := json.Unmarshal(c.PostBody(), &rawTx); err != nil {
		logger.HandlerError("SendRawTx", err)
		return err
	}

	res, err := req.Get(storage.EndpointForReq.Get() + "/sendtx/" + rawTx.Data)
	if err != nil {
		logger.HandlerError("SendRawTx", err)
		return err
	}

	if res.Response().StatusCode != 200 {
		err := errors.New("tx reverted")
		logger.HandlerError("SendRawTx", err)
		return err
	}

	r := struct {
		Result string `json:"result"`
	}{}

	if err = res.ToJSON(&r); err != nil {
		logger.HandlerError("SendRawTx", err)
		return err
	}

	result.Hash = r.Result

	if err := responses.JsonResponse(c, result); err != nil {
		return err
	}

	logger.LogRequest(time.Since(start), os.Getenv("BLOCKCHAIN"), "SendRawTx", false)

	return nil
}
