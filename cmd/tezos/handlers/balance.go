package handlers

import (
	"github.com/button-tech/utils-node-tool/types/requests"
	"github.com/button-tech/utils-node-tool/types/responses"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	routing "github.com/qiangxue/fasthttp-routing"
)

const (
	tezosURL = "https://api6.dunscan.io/v3/balance/"
)

func GetBalance(ctx *routing.Context) error {
	address := ctx.Param("address")

	balance, err := getTezosBalance(address)
	if err != nil {
		return err
	}

	if err := responses.JsonResponse(ctx, &responses.BalanceResponse{
		Balance: balance,
	}); err != nil {
		return err
	}
	return nil
}

func getTezosBalance(address string) (string, error) {
	rq := req.New()
	url := tezosURL + address

	resp, err := rq.Get(url)
	if err != nil {
		return "", errors.Wrap(err, "getTezosBalanceRequest")
	}

	b := make(requests.TezosBalance, 0)
	if err := resp.ToJSON(&b); err != nil {
		return "", errors.Wrap(err, "toJSON")
	}
	balance := b[0]

	if isNoBalance(balance) {
		return "", errors.New("balance doesn't exist")
	}

	return balance, nil
}

func isNoBalance(s string) bool {
	return s == "0"
}
