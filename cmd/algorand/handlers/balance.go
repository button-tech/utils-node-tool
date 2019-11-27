package handlers

import (
	"github.com/button-tech/utils-node-tool/types/requests"
	"github.com/button-tech/utils-node-tool/types/responses"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

// todo: node?
const urlAlgorand = "/v1/account/{address}"

func GetBalance(ctx *routing.Context) error {
	address := ctx.Param("address")

	url := urlAlgorand + address
	balance, err := getBalance(url)
	if err != nil {
		return err
	}
	return responses.JsonResponse(ctx, responses.BalanceData{Balance: balance})
}

func getBalance(url string) (int64, error) {
	rq := req.New()

	resp, err := rq.Get(url)
	if err != nil {
		return 0, errors.Wrap(err, "getBalanceAlgorandRequest")
	}

	if resp.Response().StatusCode != fasthttp.StatusOK {
		return 0, errors.Wrap(errors.New("statusCodeNotOK"), "getBalanceAlgorand")
	}

	var r requests.AlgorandBalance
	if err := resp.ToJSON(&r); err != nil {
		return 0, errors.Wrap(err, "algorandToJSON")
	}

	balance := int64(r.Amount)
	return balance, nil
}
