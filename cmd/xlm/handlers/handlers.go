package handlers

import (
	"github.com/button-tech/logger"
	"github.com/button-tech/utils-node-tool/types/responses"
	"github.com/imroc/req"
	"github.com/qiangxue/fasthttp-routing"
	"time"
)

func GetBalance(c *routing.Context) error {

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
