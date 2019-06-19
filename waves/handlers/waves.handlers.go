package handlers

import (
	"log"
	"strconv"
	"sync"

	"github.com/button-tech/utils-node-tool/shared/multiBalance"
	"github.com/button-tech/utils-node-tool/shared/responseModels"
	"github.com/imroc/req"
	"github.com/qiangxue/fasthttp-routing"
	"encoding/json"
)

func GetBalance(c *routing.Context) error {

	address := c.Param("address")

	res, err := req.Get("https://nodes.wavesplatform.com/addresses/balance/" + address)
	if err != nil {
		log.Println(err)
		return err
	}

	var data responses.BalanceData

	err = res.ToJSON(&data)
	if err != nil {
		log.Println(err)
		return err
	}

	response := new(responses.BalanceResponse)

	response.Balance = strconv.FormatInt(data.Balance, 10)

	if err := responses.JsonResponse(c, response);err != nil{
		return err
	}

	return nil
}

func GetBalances(c *routing.Context) error {

	type Request struct {
		AddressesArray []string `json:"addressesArray"`
	}

	request := new(Request)

	balances := multiBalance.New()

	if err := json.Unmarshal(c.PostBody(), &request); err != nil{
		log.Println(err)
		return err
	}

	var wg sync.WaitGroup

	for i := 0; i < len(request.AddressesArray); i++ {
		wg.Add(1)
		go multiBalance.WavesWorker(&wg, request.AddressesArray[i], balances)
	}
	wg.Wait()

	response := new(responses.BalancesResponse)

	response.Balances = balances.Result

	if err := responses.JsonResponse(c, response);err != nil{
		return err
	}

	return nil
}