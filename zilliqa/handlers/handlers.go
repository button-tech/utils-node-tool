package handlers

import (
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/button-tech/utils-node-tool/shared/responses"
	"github.com/qiangxue/fasthttp-routing"
)

func GetBalance(c *routing.Context) error {

	zilliqaAddress := c.Param("address")

	decodedAddress, err := bech32.FromBech32Addr(zilliqaAddress)
	if err != nil {
		return err
	}

	endpoint := provider.NewProvider("https://api.zilliqa.com/")

	balance := endpoint.GetBalance(decodedAddress)

	response := new(responses.BalanceResponse)

	if balance.Result == nil {
		response.Balance = "0"
		if err := responses.JsonResponse(c, response); err != nil {
			return err
		}
	}

	
	response.Balance = fmt.Sprintf("%v", balance.Result.(map[string]interface{})["balance"])

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	return nil

}
