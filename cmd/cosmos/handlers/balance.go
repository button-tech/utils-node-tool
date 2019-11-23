package handlers

import (
	"github.com/button-tech/utils-node-tool/types/requests"
	"github.com/button-tech/utils-node-tool/types/responses"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

const url = "stargate.cosmos.network/bank/balances/"

func GetBalance(ctx *routing.Context) error {
	address := ctx.Param("address")

	url := url + address
	balance, err := getBalance(url)
	if err != nil {
		return err
	}

	return responses.JsonResponse(ctx, responses.BalanceResponse{Balance: balance})
}

func getBalance(url string) (string, error) {
	rq := req.New()

	resp, err := rq.Get(url)
	if err != nil {
		return "", errors.Wrap(err, "getCosmosBalance")
	}

	if resp.Response().StatusCode != fasthttp.StatusOK {
		return "", errors.Wrap(errors.New("responseStatusNotOk"), "CosmosGetBalance")
	}

	b := make(requests.CosmosBalance, 1)
	if err = resp.ToJSON(&b); err != nil {
		return "", errors.Wrap(err, "COSMOStoJSON")
	}
	balance := b[0].Amount

	return balance, nil
}
