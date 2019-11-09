package handlers

import (
	"fmt"
	"github.com/button-tech/utils-node-tool/types/requests"
	"github.com/button-tech/utils-node-tool/types/responses"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	routing "github.com/qiangxue/fasthttp-routing"
)

const (
	xrpURL = "https://data.ripple.com/v2/accounts/%s/balances?currency=XRP"
)

func GetBalance(ctx *routing.Context) error {
	address := ctx.Param("address")

	balance, err := getXRPBalance(address)
	if err != nil {
		return err
	}
	if err := responses.JsonResponse(ctx, balance); err != nil {
		return err
	}
	return nil
}

func getXRPBalance(address string) (string, error) {
	rq := req.New()
	url := fmt.Sprintf(xrpURL, address)

	resp, err := rq.Get(url)
	if err != nil {
		return "", errors.Wrap(err, "getXRPBalanceRequest")
	}

	var b requests.XRPBalance
	if err := resp.ToJSON(&b); err != nil {
		return "", errors.Wrap(err, "toJSON")
	}

	if !checkResponseOk(b.Result) {
		return "", errors.New("No balance")
	}
	balance := b.Balances[0].Value

	return balance, nil
}

func checkResponseOk(s string) bool {
	if s == "error" {
		return false
	}
	return true
}
