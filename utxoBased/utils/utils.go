package utils

import (
	"errors"
	"fmt"
	"github.com/button-tech/utils-node-tool/shared/db"
	"github.com/button-tech/utils-node-tool/shared/responses"
	"github.com/imroc/req"
	"golang.org/x/sync/errgroup"
	"os"
	"strconv"
)

func GetBalance(address string) (string, error) {

	var reserveUrl string

	currency := os.Getenv("blockchain")

	switch currency {
	case "btc":
		reserveUrl = "/addr/" + address + "/balance"
	case "bch":
		reserveUrl = "/v1/address/details/" + address
	case "ltc":
		reserveUrl = "/api/addr/" + address + "/balance"
	}

	data := struct {
		Balance string `json:"balance"`
	}{}

	reserveData := struct {
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
			err = responseOfReserveApi.ToJSON(&reserveData)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%f", reserveData.Balance), nil

		}

		balanceFloat, err := strconv.ParseFloat(responseOfReserveApi.String(), 64)
		if err != nil {
			return "", err
		}

		balanceFloat *= 0.00000001

		return fmt.Sprintf("%f", balanceFloat), nil
	}

	err = responseOfMainApi.ToJSON(&data)
	if err != nil {
		return "", err
	}

	return data.Balance, nil
}

func GetUTXO(address string) ([]responses.UTXO, error) {

	currency := os.Getenv("blockchain")

	var endPoint string

	var requestUrl string

	if currency == "bch" {
		endPoint = "https://rest.bitbox.earth"
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

	if utxos.Response().StatusCode != 200 {
		return nil, errors.New("Bad request")
	}

	var utxoArray []responses.UTXO

	err = utxos.ToJSON(&utxoArray)
	if err != nil {
		return nil, err
	}

	return utxoArray, nil
}

func GetBalances(addresses []string) (map[string]string, error) {

	result := responses.NewBalances()

	var g errgroup.Group

	for _, address := range addresses {

		address := address

		g.Go(func() error {

			balance, err := GetBalance(address)
			if err != nil {
				return err
			}

			result.Set(address, balance)

			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result.AddressesAndBalances, nil
}
