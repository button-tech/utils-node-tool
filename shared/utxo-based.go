package shared

import (
	"errors"
	"github.com/button-tech/utils-node-tool/shared/requests"
	"github.com/button-tech/utils-node-tool/shared/responses"
	"github.com/button-tech/utils-node-tool/utils-for-endpoints/storage"
	"github.com/imroc/req"
	"log"
	"strconv"
)

func GetUtxo(address string) ([]responses.UTXO, error) {

	utxos, err := req.Get(storage.EndpointForReq.Get() + "/utxo/" + address)
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

func GetUtxoBasedBlockNumber(currency, addr string) (int64, error) {

	var (
		info requests.UtxoBasedBlocksHeight
		url  string
	)

	res, err := req.Get(addr + url)
	if err != nil || res.Response().StatusCode != 200 {
		err := DeleteEntry(currency, addr)
		if err != nil {
			return 0, err
		}
		log.Println("Status code:" + strconv.Itoa(res.Response().StatusCode))
	}

	err = res.ToJSON(&info)
	if err != nil {
		return 0, err
	}

	return info.Backend.Blocks, nil
}
