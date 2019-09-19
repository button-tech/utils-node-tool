package estorage

import (
	"errors"
	"github.com/button-tech/utils-node-tool/shared/db"
	"github.com/button-tech/utils-node-tool/shared/db/schema"
	"github.com/imroc/req"
	"github.com/onrik/ethrpc"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

type StoredEndpoints struct {
	sync.RWMutex
	Entry schema.EndpointsData
}

func (s *StoredEndpoints) Set(entry schema.EndpointsData) {
	s.Lock()
	s.Entry = entry
	s.Unlock()
}

func (s *StoredEndpoints) Get() *schema.EndpointsData {
	s.RLock()
	defer s.RUnlock()
	return &s.Entry
}

type FastestEndpoint struct {
	sync.RWMutex
	Address string
}

func (f *FastestEndpoint) Set(addr string) {
	f.Lock()
	f.Address = addr
	f.Unlock()
}

func (f *FastestEndpoint) Get() string {
	f.RLock()
	defer f.RUnlock()
	return f.Address
}

type GetFastestEndpoint func(string) string

var (
	EndpointsFromDB StoredEndpoints
	EndpointForReq  FastestEndpoint
)

func StoreEndpointsFromDB() {
	log.Println("Started storing!")
	for {
		entry, err := db.GetEntry()
		if err != nil {
			log.Fatal(err)
		}

		EndpointsFromDB.Set(*entry)

		time.Sleep(time.Minute * 10)
	}
}

func SetFastestEndpoint() {
	log.Println("Started set fastest endpoint!")

	var (
		getEndpoint GetFastestEndpoint
		mainUrl     string
	)

	switch os.Getenv("BLOCKCHAIN") {
	case "eth":
		getEndpoint = GetFastestEthBasedEndpoint
		mainUrl = os.Getenv("MAIN_API")
	case "etc":
		getEndpoint = GetFastestEthBasedEndpoint
		mainUrl = os.Getenv("MAIN_API")
	default:
		getEndpoint = GetFastestUtxoBasedEndpoint
		mainUrl = os.Getenv("MAIN_API") + "/v1/address/"
	}

	if len(os.Getenv("ADDRESS")) == 0 {
		log.Fatal(errors.New("Not set ADDRESS env!"))
	}

	for {
		result := getEndpoint(mainUrl)

		EndpointForReq.Set(result)

		time.Sleep(time.Minute * 1)
	}
}

func GetFastestUtxoBasedEndpoint(mainUrl string) string {
	currency := os.Getenv("BLOCKCHAIN")

	var endpoints []string
	switch currency {
	case "btc":
		dbEndpoints := EndpointsFromDB.Get().Addresses
		for _, j := range dbEndpoints {
			j = j + "/addr/"
			endpoints = append(endpoints, j)
		}
		endpoints = append(endpoints, mainUrl)

	case "ltc":
		dbEndpoints := EndpointsFromDB.Get().Addresses
		for _, j := range dbEndpoints {
			j = j + "/api/addr/"
			endpoints = append(endpoints, j)
		}
		endpoints = append(endpoints, mainUrl)

	case "bch":
		dbEndpoints := EndpointsFromDB.Get().Addresses
		endpoints = append(endpoints, dbEndpoints...)
		endpoints = append(endpoints, mainUrl)
	}

	fastestEndpoint := make(chan string, len(endpoints))

	for _, addr := range endpoints {
		go func(addr string) {
			res, err := req.Get(addr + os.Getenv("ADDRESS"))
			if err != nil || res.Response().StatusCode != 200 {
				return
			}

			if res.Response().StatusCode == 200 {
				fastestEndpoint <- addr
			}
		}(addr)
	}

	return <-fastestEndpoint
}

func GetFastestEthBasedEndpoint(mainUrl string) string {

	endpoints := EndpointsFromDB.Get().Addresses
	endpoints = append(endpoints, mainUrl)

	fastestEndpoint := make(chan string, len(endpoints))

	for _, e := range endpoints {
		go func(e string) {
			ethClient := ethrpc.New(e)

			_, err := ethClient.EthGetBalance(os.Getenv("ADDRESS"), "latest")
			if err != nil {
				return
			}

			fastestEndpoint <- e
		}(e)
	}

	return <-fastestEndpoint
}

func GetEndpoint(currency string) (string, error) {
	endpoints := EndpointsFromDB.Get()
	if endpoints == nil {
		return "", errors.New("Not found")
	}

	addresses := endpoints.Addresses

	rand.Seed(time.Now().UnixNano())

	result := addresses[rand.Intn(len(addresses))]

	return result, nil
}
