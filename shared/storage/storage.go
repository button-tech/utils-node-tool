package storage

import (
	"errors"
	"github.com/button-tech/utils-node-tool/shared/db"
	"github.com/button-tech/utils-node-tool/shared/db/schema"
	"github.com/imroc/req"
	"github.com/onrik/ethrpc"
	"log"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

type StoredEndpoints struct {
	sync.RWMutex
	Entry schema.EndpointsData
}

func (s *StoredEndpoints) set(entry schema.EndpointsData) {
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

func (f *FastestEndpoint) set(addr string) {
	f.Lock()
	f.Address = addr
	f.Unlock()
}

func (f *FastestEndpoint) Get() string {
	f.RLock()
	defer f.RUnlock()
	return f.Address
}

type GetFastestEndpoint func() string

var (
	EndpointsFromDB StoredEndpoints
	EndpointForReq  FastestEndpoint
)

func StoreEndpointsFromDB(startChan chan<- struct{}) {

	// For first set
	entry, err := db.GetEntry()
	if err != nil {
		log.Fatal(err)
	}

	if entry == nil {
		log.Fatal(errors.New("Something wrong with db entry!"))
	}

	EndpointsFromDB.set(*entry)

	// Send signal to start set fastest endpoint
	startChan <- struct{}{}

	log.Println("Successfully updated!")

	time.Sleep(time.Minute * 1)

	log.Println("Started storing!")

	for {
		log.Println("Trying to update...")
		entry, err := db.GetEntry()
		if err != nil || entry == nil {
			log.Println("Something wrong with entry or db!")
			time.Sleep(time.Minute * 5)
			continue
		}

		EndpointsFromDB.set(*entry)

		log.Println("Successfully updated")

		time.Sleep(time.Minute * 10)
	}
}

func SetFastestEndpoint(startChan chan struct{}) {

	<-startChan

	log.Println("Got signal from chan!")
	close(startChan)

	var (
		getEndpoint GetFastestEndpoint
	)

	switch os.Getenv("BLOCKCHAIN") {
	case "eth", "etc":
		getEndpoint = GetFastestEthBasedEndpoint
	default:
		getEndpoint = GetFastestUtxoBasedEndpoint
	}

	if len(os.Getenv("ADDRESS")) == 0 {
		log.Fatal(errors.New("Not set ADDRESS env!"))
	}

	log.Println("Started set fastest endpoint!")

	for {
		endpoint := getEndpoint()
		if len(endpoint) == 0 {
			log.Println("WARNING: All endpoints are not available now!")
			time.Sleep(time.Minute * 1)
			continue
		}

		EndpointForReq.set(endpoint)

		log.Println(EndpointForReq.Get())

		log.Println(runtime.NumGoroutine())

		time.Sleep(time.Minute * 1)
	}
}

func GetFastestUtxoBasedEndpoint() string {

	endpoints := EndpointsFromDB.Get().Addresses

	fastestEndpoint := make(chan string, len(endpoints))

	for _, addr := range endpoints {
		go func(addr string) {
			res, err := req.Get(addr + "/address/" + os.Getenv("ADDRESS"))
			if err != nil || res.Response().StatusCode != 200 {
				return
			}

			if err = res.Response().Body.Close(); err != nil {
				return
			}

			fastestEndpoint <- addr

		}(addr)
	}

	select {
	case result := <-fastestEndpoint:
		return result
	case <-time.After(time.Second * 2):
		return ""
	}
}

func GetFastestEthBasedEndpoint() string {

	endpoints := EndpointsFromDB.Get().Addresses

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

	select {
	case result := <-fastestEndpoint:
		return result
	case <-time.After(time.Second * 2):
		return ""
	}
}

func GetEndpoint() (string, error) {
	endpoints := EndpointsFromDB.Get()
	if endpoints == nil {
		return "", errors.New("Not found")
	}

	addresses := endpoints.Addresses

	rand.Seed(time.Now().UnixNano())

	result := addresses[rand.Intn(len(addresses))]

	return result, nil
}
