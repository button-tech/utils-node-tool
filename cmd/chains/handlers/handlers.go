package handlers

import (
	"fmt"
	"github.com/button-tech/logger"
	"github.com/button-tech/utils-node-tool/types/requests"
	"github.com/button-tech/utils-node-tool/types/responses"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
	"time"
)

const (
	urlBnBTestnet = "https://testnet-dex.binance.org/api/v1/account/"
	urlBnBMainnet = "https://dex.binance.org/api/v1/account/"
)

func GetCosmosBalance(ctx *routing.Context) error {
	address := ctx.Param("address")

	url := "stargate.cosmos.network/bank/balances/" + address

	rq := req.New()

	resp, err := rq.Get(url)
	if err != nil {
		return errors.Wrap(err, "getCosmosBalance")
	}

	if resp.Response().StatusCode != fasthttp.StatusOK {
		return errors.Wrap(errors.New("responseStatusNotOk"), "CosmosGetBalance")
	}

	b := make(requests.CosmosBalance, 1)
	if err = resp.ToJSON(&b); err != nil {
		return errors.Wrap(err, "COSMOStoJSON")
	}
	balance := b[0].Amount

	return responses.JsonResponse(ctx, responses.BalanceResponse{Balance: balance})
}

func GetWavesBalance(c *routing.Context) error {

	start := time.Now()

	address := c.Param("address")

	res, err := req.Get("https://nodes.wavesplatform.com/addresses/balance/" + address)
	if err != nil {
		return err
	}

	var data responses.BalanceData

	err = res.ToJSON(&data)
	if err != nil {
		return err
	}

	response := new(responses.BalanceResponse)

	response.Balance = strconv.FormatInt(data.Balance, 10)

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	logger.LogRequest(time.Since(start), "WAVES", "GetBalance", false)

	return nil
}

func GetXlmBalance(c *routing.Context) error {

	start := time.Now()

	type StellarBalance struct {
		Balances []struct {
			Balance             string `json:"balance"`
			Buying_liabilities  string `json:"buying_liabilities"`
			Selling_liabilities string `json:"selling_liabilities"`
			Asset_type          string `json:"asset_type"`
		}
	}

	var balance StellarBalance

	balanceReq, err := req.Get("https://horizon.stellar.org/accounts/" + c.Param("address"))
	if err != nil {
		return err
	}

	err = balanceReq.ToJSON(&balance)
	if err != nil {
		return err
	}

	var stellarBalanceString string

	for _, j := range balance.Balances {
		if j.Asset_type == "native" {
			stellarBalanceString = j.Balance
		}
	}

	if stellarBalanceString == "" {
		stellarBalanceString = "0"
	}

	response := new(responses.BalanceResponse)

	response.Balance = stellarBalanceString

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	logger.LogRequest(time.Since(start), "XLM", "GetBalance", false)

	return nil
}

func GetXrpBalance(ctx *routing.Context) error {
	address := ctx.Param("address")

	rq := req.New()
	url := fmt.Sprintf("https://data.ripple.com/v2/accounts/%s/balances?currency=XRP", address)

	resp, err := rq.Get(url)
	if err != nil {
		return errors.Wrap(err, "getXRPBalanceRequest")
	}

	var b requests.XRPBalance
	if err := resp.ToJSON(&b); err != nil {
		return errors.Wrap(err, "toJSON")
	}

	if b.Result == "error" {
		return errors.New("No balance")
	}
	balance := b.Balances[0].Value

	return responses.JsonResponse(ctx, &responses.BalanceResponse{Balance: balance})
}

func GetTezosBalance(ctx *routing.Context) error {
	address := ctx.Param("address")

	rq := req.New()
	url := "https://api6.dunscan.io/v3/balance/" + address

	resp, err := rq.Get(url)
	if err != nil {
		return errors.Wrap(err, "getTezosBalanceRequest")
	}

	b := make(requests.TezosBalance, 1)
	if err := resp.ToJSON(&b); err != nil {
		return errors.Wrap(err, "toJSON")
	}

	balance := b[0]

	if balance == "0" {
		return errors.New("balance doesn't exist")
	}

	return responses.JsonResponse(ctx, &responses.BalanceResponse{Balance: balance})
}

func GetBnbBalance(ctx *routing.Context) error {
	address := ctx.Param("address")

	var url string
	if strings.HasPrefix(address, "bnb") {
		url = urlBnBMainnet
	} else {
		url = urlBnBTestnet
	}

	url += address
	rq := req.New()
	resp, err := rq.Get(url)
	if err != nil {
		return errors.Wrap(err, "getBnBBalanceRequest")
	}

	var b requests.BnBBalance
	if err := resp.ToJSON(&b); err != nil {
		return errors.Wrap(err, "toJSON")
	}

	var balance string
	for _, v := range b.Balances {
		if v.Symbol == "BNB" {
			balance = v.Free
		}
	}

	if balance == "" {
		return errors.New("No BnB balance")
	}

	return responses.JsonResponse(ctx, &responses.BalanceResponse{Balance: balance})
}

func GetAlgorandBalance(ctx *routing.Context) error {
	address := ctx.Param("address")

	// todo: node?
	url := "/v1/account/{address}" + address
	rq := req.New()

	resp, err := rq.Get(url)
	if err != nil {
		return errors.Wrap(err, "getBalanceAlgorandRequest")
	}

	if resp.Response().StatusCode != fasthttp.StatusOK {
		return errors.Wrap(errors.New("statusCodeNotOK"), "getBalanceAlgorand")
	}

	var r requests.AlgorandBalance
	if err := resp.ToJSON(&r); err != nil {
		return errors.Wrap(err, "algorandToJSON")
	}

	balance := int64(r.Amount)
	return responses.JsonResponse(ctx, responses.BalanceData{Balance: balance})
}
