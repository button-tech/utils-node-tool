package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/button-tech/utils-node-tool/shared"
	"github.com/button-tech/utils-node-tool/shared/db"
	"golang.org/x/sync/errgroup"
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
		blockDifference int64 = 10
		result          Result
	)

	switch currency {
	case "eth","etc":
		getBlockNumber = shared.GetEthBasedBlockNumber
	default:
		getBlockNumber = shared.GetUtxoBasedBlockNumber
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
			log.Println("BlockChainHeight:" + strconv.Itoa(int(j.BlockChainHeight)))
			log.Println("Sync now:" + strconv.Itoa(int(maxNumber)))
		}
	}

	fmt.Println("All " + currency + " nodes checked! Alive nodes count - " + strconv.Itoa(len(result.NodesInfo)))
	return nil
}

func main() {

	log.Println("Start!")

	for {

		entries, err := db.GetAll()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		var g errgroup.Group

		for _, j := range entries {
			j := j
			g.Go(func() error {
				return SyncCheck(j.Currency, j.Addresses)
			})

		}
		if err := g.Wait(); err != nil {
			log.Println(err)
			os.Exit(1)
		}

		time.Sleep(time.Minute * 10)
	}

}
