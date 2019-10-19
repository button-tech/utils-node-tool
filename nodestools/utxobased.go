package nodestools

import (
	"errors"
	"github.com/button-tech/utils-node-tool/requests"
	"github.com/button-tech/utils-node-tool/responses"
	"github.com/button-tech/utils-node-tool/storage"
	"github.com/imroc/req"
	"golang.org/x/sync/errgroup"
)

func GetUtxo(address string) ([]responses.UTXO, error) {

	endpoint := storage.EndpointForReq.Get()

	utxos, err := req.Get(endpoint + "/utxo/" + address)
	if err != nil {
		return nil, err
	}

	if utxos.Response().StatusCode != 200 {
		return nil, errors.New("Bad request")
	}

	var (
		utxoArray []responses.UTXO
		g         errgroup.Group
	)

	err = utxos.ToJSON(&utxoArray)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(utxoArray); i++ {
		i := i
		g.Go(func() error {

			var tx requests.UtxoBasedTxOutputs

			res, err := req.Get(endpoint + "/tx/" + utxoArray[i].Txid)
			if err != nil {
				return err
			}

			if res.Response().StatusCode != 200 {
				return errors.New("Bad request!")
			}

			err = res.ToJSON(&tx)
			if err != nil {
				return err
			}

			for _, el := range tx.Vout {
				if Contains(el.ScriptPubKey.Addresses, address) {
					utxoArray[i].ScriptPubKey = el.ScriptPubKey.Hex
					utxoArray[i].Address = address
				}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return utxoArray, nil
}

func GetUtxoBasedBlockNumber(addr string) (int64, error) {

	var (
		info requests.UtxoBasedBlocksHeight
		url  string
	)

	res, err := req.Get(addr + url)
	if err != nil {
		return 0, err
	}

	if res.Response().StatusCode != 200 {
		return 0, errors.New("Bad request")
	}

	err = res.ToJSON(&info)
	if err != nil {
		return 0, err
	}

	return info.Backend.Blocks, nil
}
