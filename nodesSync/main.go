package main

import (
	"fmt"
	"github.com/button-tech/utils-node-tool/shared"
	"github.com/button-tech/utils-node-tool/shared/db"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"time"
)

type SyncStatus struct {
	BlockChainHeight int64
	Address          string
}

type Req func(address, currency string) (int64, error)

func SyncCheck(currency string, addresses []string) error {

	var results []SyncStatus

	var numbers []int64

	var getBlockNumber Req

	if currency == "btc" || currency == "ltc" {
		getBlockNumber = shared.GetUtxoBasedBlockNumber
	} else {
		getBlockNumber = shared.GetEthBasedBlockNumber
	}

	var g errgroup.Group

	for _, addr := range addresses {
		addr := addr
		g.Go(func() error {
			blockNumber, err := getBlockNumber(currency, addr)
			if err != nil {
				return err
			}
			results = append(results, SyncStatus{blockNumber, addr})
			numbers = append(numbers, blockNumber)

			return nil
		})
	}
	if err := g.Wait(); err != nil {
		log.Println(err)
		return err
	}

	maxNumber := shared.Max(numbers)

	fmt.Println(results)

	for _, j := range results {
		if j.BlockChainHeight < maxNumber-3 {
			err := shared.DeleteEntry(currency, j.Address)
			if err != nil {
				return err
			}
		}
	}

	fmt.Println("All " + currency + " nodes checked!")

	return nil
}

func main() {

	for {

		entrys, err := db.GetAll()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		var g errgroup.Group

		for _, j := range entrys {
			j := j
			g.Go(func() error {
				err := SyncCheck(j.Currency, j.Addresses)
				if err != nil {
					return err
				}
				return nil
			})

		}
		if err := g.Wait(); err != nil {
			log.Println(err)
			os.Exit(1)
		}

		time.Sleep(time.Second * 10)
	}

}
