package handlers

import (
	"strings"

	"github.com/button-tech/utils-node-tool/types/requests"
	"github.com/button-tech/utils-node-tool/types/responses"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	routing "github.com/qiangxue/fasthttp-routing"
)

const (
	urlBnBTestnet = "https://testnet-dex.binance.org/api/v1/account/"
	urlBnBMainnet = "https://dex.binance.org/api/v1/account/"
)

func GetBalance(ctx *routing.Context) error {
	address := ctx.Param("address")

	var url string
	if mainnet(address) {
		url = urlBnBMainnet
	} else {
		url = urlBnBTestnet
	}

	url += address
	balance, err := getBnBBalance(url)
	if err != nil {
		return err
	}

	return responses.JsonResponse(ctx, &responses.BalanceResponse{Balance: balance})
}

func getBnBBalance(url string) (string, error) {
	rq := req.New()
	resp, err := rq.Get(url)
	if err != nil {
		return "", errors.Wrap(err, "getBnBBalanceRequest")
	}

	var b requests.BnBBalance
	if err := resp.ToJSON(&b); err != nil {
		return "", errors.Wrap(err, "toJSON")
	}

	return findBnB(&b)
}

func mainnet(address string) bool {
	return strings.HasPrefix(address, "bnb")
}

func findBnB(b *requests.BnBBalance) (string, error) {
	var balance string
	for _, v := range b.Balances {
		if v.Symbol == "BNB" {
			balance = v.Free
			return balance, nil
		}
	}
	return "", errors.New("No BnB balance")
}
