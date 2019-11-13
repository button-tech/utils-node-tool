package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/button-tech/utils-node-tool/types/requests"
	"github.com/button-tech/utils-node-tool/types/responses"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	routing "github.com/qiangxue/fasthttp-routing"
)

const (
	xrpTestnet = "altnet"
	xrpDevnet  = "devnet"

	xrpSendTxURL = "https://s.%s.rippletest.net:51234/submit"
)

func SubmitTransaction(ctx *routing.Context) error {
	var txData requests.XRPTxData
	if err := json.Unmarshal(ctx.PostBody(), &txData); err != nil {
		return err
	}

	url := setupURL(txData.Devnet)
	submitted, err := submitTx(url, xrpTxToSubmit(txData.TxBlob))
	if err != nil {
		return err
	}

	if err := checkSubmitTxStatus(submitted); err != nil {
		return err
	}

	txHash := submitted.Result.TxJSON.Hash
	if err := responses.JsonResponse(ctx, responses.XRPTxHash{TxHash: txHash}); err != nil {
		return err
	}
	return nil
}

func setupURL(devnet bool) string {
	if devnet {
		return fmt.Sprintf(xrpSendTxURL, xrpDevnet)
	}
	return fmt.Sprintf(xrpSendTxURL, xrpTestnet)
}

const submitMethod = "submit"

func xrpTxToSubmit(txBlob string) requests.XRPTxToSubmit {
	return requests.XRPTxToSubmit{
		Method: submitMethod,
		Params: []struct {
			TxBlob string `json:"tx_blob"`
		}{
			{
				TxBlob: txBlob,
			},
		},
	}
}

func submitTx(url string, data requests.XRPTxToSubmit) (*requests.XRPSentTxInfo, error) {
	rq := req.New()

	resp, err := rq.Post(url, req.BodyJSON(&data))
	if err != nil {
		return nil, errors.Wrap(err, "submitTxRequest")
	}

	if resp.Response().StatusCode != 200 {
		return nil, errors.Wrap(errors.New("StatusCodeNotOk"), "Request to Ripple")
	}

	var info requests.XRPSentTxInfo
	if err := resp.ToJSON(&info); err != nil {
		return nil, errors.Wrap(err, "XRPtoJSON")
	}
	return &info, nil
}

func checkSubmitTxStatus(info *requests.XRPSentTxInfo) error {
	if info.Result.Status == "error" {
		return errors.Wrap(errors.New("ResponseStatusError"), "RippleAPI")
	}
	return nil
}
