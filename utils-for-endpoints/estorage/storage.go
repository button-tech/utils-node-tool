package estorage

import (
	"errors"
	"github.com/button-tech/utils-node-tool/shared/db"
	"github.com/button-tech/utils-node-tool/shared/db/schema"
	"log"
	"math/rand"
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

var (
	EndpointsFromDB StoredEndpoints
	//EndpointsForReq FastestEndpoint
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

//func SetFastestEndpoints(){
//	log.Println("Started set fastest endpoints!")
//	for {
//		time.Sleep(time.Minute * 10)
//	}
//}
//
//func GetFastestUtxoBasedEndpoint(currency string, endpoints []string) string {
//	currency := os.Getenv("BLOCKCHAIN")
//
//	mainApi := os.Getenv("MAIN_API")
//
//	mainUrl := mainApi + "/v1/address/" + address
//
//	switch currency {
//	case "btc":
//		dbEndpoints := estorage.EndpointsFromDB.BtcEndpoints.Addresses
//		for _, j := range dbEndpoints {
//			j = j + "/addr/" + address
//			endpoints = append(endpoints, j)
//		}
//		endpoints = append(endpoints, mainUrl)
//	case "ltc":
//		dbEndpoints := estorage.EndpointsFromDB.LtcEndpoints.Addresses
//		for _, j := range dbEndpoints {
//			j = j + "/api/addr/" + address
//			endpoints = append(endpoints, j)
//		}
//		endpoints = append(endpoints, mainUrl)
//	case "bch":
//		endpoints = append(endpoints, mainUrl)
//		endpoints = append(endpoints, "https://rest.bitbox.earth/v1/address/details/"+address)
//	}
//
//	fastestEndpoint := make(chan string, len(endpoints))
//
//	for _, addr := range endpoints {
//		go func(addr string) {
//			res, err := req.Get(addr + "/" + os.Getenv(strings.ToUpper(currency) + "_ADDRESS") + "/")
//			if err != nil || res.Response().StatusCode != 200 {
//				return
//			}
//
//			if res.Response().StatusCode == 200{
//				fastestEndpoint <- addr
//			}
//		}(addr)
//	}
//
//	select {
//	case result := <-fastestEndpoint:
//		return result
//	}
//}
//
//func GetFastestEthBasedEndpoint(endpoints []string){
//	balanceChan := make(chan string, len(endpoints))
//
//	for _, e := range endpoints {
//		go func(e string) {
//			ethClient := ethrpc.New(e)
//			res, err := ethClient.EthGetBalance(os.Getenv("ETH_ADDRESS"), "latest")
//			if err != nil {
//				return
//			}
//
//			balanceChan <- res.String()
//		}(e)
//	}
//
//	select {
//	case result := <-balanceChan:
//		return result, nil
//	case <-time.After(2 * time.Second):
//		return "", errors.New("Bad request")
//	}
//}
//
//func EtherBalanceReq(endpoints []string, address string) (string, error) {
//
//}

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
