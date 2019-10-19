package handlers

import (
	"github.com/button-tech/utils-node-tool/types/responses"
	"github.com/imroc/req"
	"github.com/qiangxue/fasthttp-routing"
	"strconv"
)

func GetBalance(c *routing.Context) error {

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

	return nil
}
