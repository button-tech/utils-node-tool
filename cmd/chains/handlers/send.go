package handlers

import (
	"encoding/json"
	"errors"
	"github.com/button-tech/utils-node-tool/types/requests"
	"github.com/button-tech/utils-node-tool/types/responses"
	"github.com/imroc/req"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/stellar/go/clients/horizon"
	"strings"
)

type wavesDataTxToSubmit struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}

func SendWavesRawTx(c *routing.Context) error {
	var (
		rawTx  requests.RawTransaction
		result responses.TransactionResult
	)

	if err := json.Unmarshal(c.PostBody(), &rawTx); err != nil {
		return err
	}

	url := "https://nodes.wavesplatform.com/transactions/broadcast"

	payload := strings.NewReader(rawTx.Data)
	res, err := req.Post(url, req.Header{"Content-Type": "application/json"}, payload)
	if err != nil {
		return err
	}

	var r wavesDataTxToSubmit
	if err = res.ToJSON(&r); err != nil {
		return err
	}

	if len(r.Message) != 0 {
		return errors.New(r.Message)
	}

	result.Hash = r.ID

	if err := responses.JsonResponse(c, result); err != nil {
		return err
	}

	return nil
}

func SendXlmRawTx(c *routing.Context) error {
	var (
		rawTx  requests.RawTransaction
		result responses.TransactionResult
	)

	if err := json.Unmarshal(c.PostBody(), &rawTx); err != nil {
		return err
	}

	resp, err := horizon.DefaultPublicNetClient.SubmitTransaction(rawTx.Data)
	if err != nil {
		return err
	}

	result.Hash = resp.Hash

	if err := responses.JsonResponse(c, result); err != nil {
		return err
	}

	return nil
}
