package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/button-tech/utils-node-tool/db"
	"github.com/button-tech/utils-node-tool/nodetools"
	"github.com/imroc/req"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type nodeInfo struct {
	blockChainHeight int64
	endPoint         string
}

type result struct {
	sync.Mutex
	nodesInfo    []nodeInfo
	blockNumbers []int64
	badEndpoints []string
}

func (s *result) addToBlockNumbers(blockNumber int64) {
	s.Lock()
	s.blockNumbers = append(s.blockNumbers, blockNumber)
	s.Unlock()
}

func (s *result) addToNodesInfo(address string, blockNumber int64) {
	s.Lock()
	s.nodesInfo = append(s.nodesInfo, nodeInfo{blockNumber, address})
	s.Unlock()
}

func (s *result) addToBadEndpoints(address string) {
	s.Lock()
	s.badEndpoints = append(s.badEndpoints, address)
	s.Unlock()
}

func (s *result) clearBadEndpoints() {
	s.Lock()
	s.badEndpoints = s.badEndpoints[:0]
	s.Unlock()
}

type reqBlockNumber func(address string) (int64, error)

func syncCheck(currency string, addresses []string) error {

	var (
		getBlockNumber  reqBlockNumber
		blockDifference int64 = 10
		result          result
	)

	switch currency {
	case "eth", "etc":
		getBlockNumber = nodetools.GetEthBasedBlockNumber
	default:
		// Check only eth based endpoints
		return nil
	}

	var (
		g errgroup.Group
	)

	for _, addr := range addresses {
		addr := addr
		g.Go(func() error {
			blockNumber, err := getBlockNumber(addr)
			if err != nil {
				result.addToBadEndpoints(addr)
				return nil
			}

			result.addToBlockNumbers(blockNumber)
			result.addToNodesInfo(addr, blockNumber)

			return nil
		})
	}
	if err := g.Wait(); err != nil {
		log.Println(err)
		return err
	}

	if len(result.badEndpoints) > 0 {
		err := DeleteEntries(result.badEndpoints, currency)
		if err != nil {
			return err
		}
	}

	result.clearBadEndpoints()

	maxNumber := max(result.blockNumbers)

	for _, j := range result.nodesInfo {
		if j.blockChainHeight < maxNumber-blockDifference {
			result.addToBadEndpoints(j.endPoint)
			log.Println("BlockChainHeight:" + strconv.Itoa(int(j.blockChainHeight)))
			log.Println("Sync now:" + strconv.Itoa(int(maxNumber)))
		}
	}

	if len(result.badEndpoints) > 0 {
		err := DeleteEntries(result.badEndpoints, currency)
		if err != nil {
			return err
		}
	}

	fmt.Println("All " + currency + " nodes checked! Alive nodes count - " + strconv.Itoa(len(result.nodesInfo)-len(result.badEndpoints)))
	return nil
}

func DeleteEntries(addresses []string, currency string) error {
	for _, v := range addresses {
		err := db.DeleteEntry(currency, v)
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
			time.Sleep(time.Minute * 1)
			continue
		}

		var g errgroup.Group

		for _, j := range entries {
			j := j
			g.Go(func() error {
				return syncCheck(j.Currency, j.Addresses)
			})

		}
		if err := g.Wait(); err != nil {
			log.Println(err)
		}

		time.Sleep(time.Minute * 10)
	}

}

func max(array []int64) int64 {
	var max = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
	}
	return max
}
