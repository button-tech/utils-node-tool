package handlers

import (
	"github.com/button-tech/utils-node-tool/shared/responseModels"
	"github.com/imroc/req"
	"log"
	"github.com/qiangxue/fasthttp-routing"
)

func GetBalance(c *routing.Context) error {

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
		log.Println(err)
		return err
	}

	err = balanceReq.ToJSON(&balance)
	if err != nil {
		log.Println(err)
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

	if err := responses.JsonResponse(c, response);err != nil{
		return err
	}

	return nil
}
