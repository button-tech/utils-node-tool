package handlers

import (
	"encoding/json"
	"github.com/button-tech/utils-node-tool/nodetools"
	b "github.com/button-tech/utils-node-tool/nodetools"
	"github.com/button-tech/utils-node-tool/types/requests"
	"github.com/button-tech/utils-node-tool/types/responses"
	"github.com/qiangxue/fasthttp-routing"
)

func GetBalance(c *routing.Context) error {

	address := c.Param("address")

	balance, err := b.GetUtxoBasedBalance(address)
	if err != nil {
		return err
	}

	response := responses.BalanceResponse{Balance: balance}

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	return nil
}

func GetUtxo(c *routing.Context) error {

	address := c.Param("address")

	utxoArray, err := nodetools.GetUtxo(address)
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

	response, err := b.GetUtxoBasedBalancesByList(request.Addresses)
	if err != nil {
		return err
	}

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	return nil
}
