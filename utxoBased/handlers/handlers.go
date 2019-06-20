package handlers

import (
	"encoding/json"
	"github.com/button-tech/utils-node-tool/shared/requests"
	"github.com/button-tech/utils-node-tool/shared/responses"
	"github.com/button-tech/utils-node-tool/utxoBased/utils"
	"github.com/qiangxue/fasthttp-routing"
)

func GetBalance(c *routing.Context) error {

	address := c.Param("address")

	balance, err := utils.GetBalance(address)
	if err != nil {
		return err
	}

	response := responses.BalanceResponse{Balance: balance}

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	return nil
}

func GetUTXO(c *routing.Context) error {

	address := c.Param("address")

	utxoArray, err := utils.GetUTXO(address)
	if err != nil {
		return err
	}

	response := responses.UTXOResponse{Utxo: utxoArray}

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	return nil
}

func GetBalances(c *routing.Context) error {

	request := new(requests.BalancesRequest)

	if err := json.Unmarshal(c.PostBody(), &request); err != nil {
		return err
	}

	response, err := utils.GetBalances(request.Addresses)
	if err != nil {
		return err
	}

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	return nil
}
