package nodetools

import (
	"errors"
	"github.com/button-tech/utils-node-tool/nodetools/storage"
	"github.com/button-tech/utils-node-tool/types/requests"
	"github.com/button-tech/utils-node-tool/types/responses"
	"github.com/imroc/req"
	"golang.org/x/sync/errgroup"
	"sync"
	"time"
)

// Balances(by addresses list)
type Balances struct {
	sync.RWMutex
	AddressesAndBalances map[string]string
}

func NewBalances() *Balances {
	return &Balances{
		AddressesAndBalances: make(map[string]string),
	}
}

func (ds *Balances) SetBalances(address, balance string) {
	ds.Lock()
	ds.AddressesAndBalances[address] = balance
	ds.Unlock()
}

func GetUtxoBasedBalancesByList(addresses []string) (map[string]string, error) {

	result := NewBalances()

	var g errgroup.Group

	for _, address := range addresses {

		address := address

		g.Go(func() error {

			balance, err := GetUtxoBasedBalance(address)
			if err == nil {
				result.SetBalances(address, balance)
			}

			return err
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result.AddressesAndBalances, nil
}

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

// UTXO based blockchain - BTC, LTC, BCH
func GetUtxoBasedBalance(address string) (string, error) {

	var (
		s      requests.UtxoBasedBalance
		result string
	)

	res, err := req.Get(storage.EndpointForReq.Get() + "/address/" + address)
	if err != nil || res.Response().StatusCode != 200 {
		result, err = UtxoBasedBalanceReq(address)
		if err != nil {
			return "", err
		}

		return result, nil
	}

	err = res.ToJSON(&s)
	if err != nil {
		return "", err
	}

	return s.Balance, nil
}

func UtxoBasedBalanceReq(address string) (string, error) {

	endpoints := storage.EndpointsFromDB.Get().Addresses

	balanceChan := make(chan string, len(endpoints))

	for _, addr := range endpoints {
		go func(addr string) {

			var s requests.UtxoBasedBalance

			res, err := req.Get(addr + "/address/" + address)
			if err != nil || res.Response().StatusCode != 200 {
				return
			}

			err = res.ToJSON(&s)
			if err != nil {
				return
			}

			balanceChan <- s.Balance

		}(addr)
	}

	select {
	case result := <-balanceChan:
		return result, nil
	case <-time.After(2 * time.Second):
		return "", errors.New("Bad request")
	}
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

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
