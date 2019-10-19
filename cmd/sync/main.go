package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/button-tech/utils-node-tool/db"
	"github.com/button-tech/utils-node-tool/nodetools"
	"github.com/imroc/req"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type NodeInfo struct {
	BlockChainHeight int64
	EndPoint         string
}

type Result struct {
	sync.Mutex
	NodesInfo    []NodeInfo
	BlockNumbers []int64
	BadEndpoints []string
}

func (s *Result) AddToBlockNumbers(address string, blockNumber int64) {
	s.Lock()
	s.BlockNumbers = append(s.BlockNumbers, blockNumber)
	s.Unlock()
}

func (s *Result) AddToNodesInfo(address string, blockNumber int64) {
	s.Lock()
	s.NodesInfo = append(s.NodesInfo, NodeInfo{blockNumber, address})
	s.Unlock()
}

func (s *Result) AddToBadEndpoints(address string) {
	s.Lock()
	s.BadEndpoints = append(s.BadEndpoints, address)
	s.Unlock()
}

func (s *Result) ClearBadEndpoints() {
	s.Lock()
	s.BadEndpoints = s.BadEndpoints[:0]
	s.Unlock()
}

type Req func(address string) (int64, error)

func SyncCheck(currency string, addresses []string) error {

	var (
		getBlockNumber  Req
		blockDifference int64 = 10
		result          Result
	)

	switch currency {
	case "eth", "etc":
		getBlockNumber = nodetools.GetEthBasedBlockNumber
	default:
		// Check only eth based endpoints
		return nil
		//getBlockNumber = shared.GetUtxoBasedBlockNumber
	}

	var (
		g errgroup.Group
	)

	for _, addr := range addresses {
		addr := addr
		g.Go(func() error {
			blockNumber, err := getBlockNumber(addr)
			if err != nil {
				result.AddToBadEndpoints(addr)
				return nil
			}

			result.AddToBlockNumbers(addr, blockNumber)
			result.AddToNodesInfo(addr, blockNumber)

			return nil
		})
	}
	if err := g.Wait(); err != nil {
		log.Println(err)
		return err
	}

	if len(result.BadEndpoints) > 0 {
		err := DeleteEntries(result.BadEndpoints, currency)
		if err != nil {
			return err
		}
	}

	result.ClearBadEndpoints()

	maxNumber := nodetools.Max(result.BlockNumbers)

	for _, j := range result.NodesInfo {
		if j.BlockChainHeight < maxNumber-blockDifference {
			result.AddToBadEndpoints(j.EndPoint)
			log.Println("BlockChainHeight:" + strconv.Itoa(int(j.BlockChainHeight)))
			log.Println("Sync now:" + strconv.Itoa(int(maxNumber)))
		}
	}

	if len(result.BadEndpoints) > 0 {
		err := DeleteEntries(result.BadEndpoints, currency)
		if err != nil {
			return err
		}
	}

	fmt.Println("All " + currency + " nodes checked! Alive nodes count - " + strconv.Itoa(len(result.NodesInfo)-len(result.BadEndpoints)))
	return nil
}

func DeleteEntries(addresses []string, currency string) error {
	for _, v := range addresses {
		err := nodetools.DeleteEntry(currency, v)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func main() {

	req.SetClient(&http.Client{})

	log.Println("Started checking!")

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
