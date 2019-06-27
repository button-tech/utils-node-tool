package main

import (
	"fmt"
	"github.com/button-tech/utils-node-tool/shared"
	"github.com/button-tech/utils-node-tool/shared/db"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type NodeInfo struct {
	BlockChainHeight int64
	EndPoint         string
}

type Result struct {
	sync.Mutex
	NodesInfo    []NodeInfo
	BlockNumbers []int64
}

func (s *Result) Add(address string, blockNumber int64) {
	s.Lock()
	s.NodesInfo = append(s.NodesInfo, NodeInfo{blockNumber, address})
	s.BlockNumbers = append(s.BlockNumbers, blockNumber)
	s.Unlock()
}

type Req func(address, currency string) (int64, error)

func SyncCheck(currency string, addresses []string) error {

	var (
		getBlockNumber  Req
		blockDifference int64
		result          Result
	)

	if currency == "btc" || currency == "ltc" {
		blockDifference = 1
		getBlockNumber = shared.GetUtxoBasedBlockNumber
	} else {
		blockDifference = 5
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

			result.Add(addr, blockNumber)

			return nil
		})
	}
	if err := g.Wait(); err != nil {
		log.Println(err)
		return err
	}

	maxNumber := shared.Max(result.BlockNumbers)

	for _, j := range result.NodesInfo {
		if j.BlockChainHeight < maxNumber-blockDifference {
			err := shared.DeleteEntry(currency, j.EndPoint)
			if err != nil {
				return err
			}
		}
	}

	fmt.Println("All " + currency + " nodes checked! Alive nodes count - " + strconv.Itoa(len(result.NodesInfo)))

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

		time.Sleep(time.Second * 30)
	}

}
