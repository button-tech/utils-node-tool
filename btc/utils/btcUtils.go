package utils

import (
	"errors"
	"fmt"
	"github.com/button-tech/utils-node-tool/shared/db"
	"github.com/button-tech/utils-node-tool/shared/responseModels"
	"github.com/imroc/req"
	"os"
	"strconv"
)

func GetBtcBlockChainBalance(address string) (string, error) {

	var reserveUrl string

	currency := os.Getenv("blockChain")

	switch currency {
	case "btc":
		reserveUrl = "/addr/" + address + "/balance"
	case "bch":
		reserveUrl = "/v1/address/details/" + address
	case "ltc":
		reserveUrl = "/api/addr/" + address + "/balance"
	}

	btc := struct {
		Balance string `json:"balance"`
	}{}

	reserveBtc := struct {
		Balance float64 `json:"balance"`
	}{}

	responseOfMainApi, err := req.Get(os.Getenv("main-api") + "/v1/address/" + address)

	if err != nil || responseOfMainApi.Response().StatusCode != 200 {
		endpoint, err := db.GetEndpoint(currency)
		if err != nil {
			return "", err
		}

		responseOfReserveApi, err := req.Get(endpoint + reserveUrl)
		if err != nil {
			return "", err
		}

		if responseOfReserveApi.Response().StatusCode != 200 {
			return "", errors.New("Bad request")
		}

		if currency == "bch" {
			err = responseOfReserveApi.ToJSON(&reserveBtc)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%f", reserveBtc.Balance), nil

		}

		balanceFloat, err := strconv.ParseFloat(responseOfReserveApi.String(), 64)
		if err != nil {
			return "", err
		}

		balanceFloat *= 0.00000001

		return fmt.Sprintf("%f", balanceFloat), nil
	}

	err = responseOfMainApi.ToJSON(&btc)
	if err != nil {
		return "", err
	}

	return btc.Balance, nil
}

func GetUTXO(address string) ([]responses.UTXO, error) {

	currency := os.Getenv("blockChain")

	var endPoint string

	var requestUrl string

	if currency == "bch" {
		endPoint = os.Getenv("reserve-api")
	} else {
		nodeFromDB, err := db.GetEndpoint(currency)
		if err != nil {
			return nil, err
		}
		endPoint = nodeFromDB
	}

	switch currency {
	case "btc":
		requestUrl = endPoint + "/addr/" + address + "/utxo"
	case "bch":
		requestUrl = endPoint + "/v1/address/utxo/" + address
	case "ltc":
		requestUrl = endPoint + "/api/addr/" + address + "/utxo"
	}

	utxos, err := req.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	var utxoArray []responses.UTXO

	err = utxos.ToJSON(&utxoArray)
	if err != nil {
		return nil, err
	}

	return utxoArray, nil
}
