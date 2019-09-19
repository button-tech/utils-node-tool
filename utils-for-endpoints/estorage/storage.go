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
	BtcEndpoints,
	LtcEndpoints,
	EthEndpoints,
	EtcEndpoints schema.Endpoints
}

func (s *StoredEndpoints) Add(entry schema.Endpoints) {
	s.Lock()
	switch entry.Currency {
	case "btc":
		s.BtcEndpoints = entry
	case "ltc":
		s.LtcEndpoints = entry
	case "eth":
		s.EthEndpoints = entry
	case "etc":
		s.EtcEndpoints = entry
	}
	s.Unlock()
}

func (s *StoredEndpoints) GetByCurrency(currency string) *schema.Endpoints {
	s.RLock()
	defer s.RUnlock()
	switch currency {
	case "btc":
		return &s.BtcEndpoints
	case "ltc":
		return &s.LtcEndpoints
	case "eth":
		return &s.EthEndpoints
	case "etc":
		return &s.EtcEndpoints
	default:
		return nil
	}
}

func (s *StoredEndpoints) GetListOfAllEndpoints() [4]*schema.Endpoints {
	s.RLock()
	defer s.RUnlock()
	return [4]*schema.Endpoints{
		&s.EthEndpoints,
		&s.EtcEndpoints,
		&s.BtcEndpoints,
		&s.LtcEndpoints,
	}
}

type FastestEndpoints struct {
	sync.RWMutex
	BtcEndpoint,
	LtcEndpoint,
	EthEndpoint,
	EtcEndpoint string
}

func (f *FastestEndpoints) Add(c, e string) {
	f.Lock()
	switch c {
	case "btc":
		f.BtcEndpoint = e
	case "ltc":
		f.LtcEndpoint = e
	case "eth":
		f.EthEndpoint = e
	case "etc":
		f.EtcEndpoint = e
	}
	f.Unlock()
}

func (f *FastestEndpoints) Get(c string) string {
	f.RLock()
	defer f.RUnlock()
	switch c {
	case "btc":
		return f.BtcEndpoint
	case "ltc":
		return f.LtcEndpoint
	case "eth":
		return f.EthEndpoint
	case "etc":
		return f.EtcEndpoint
	default:
		return ""
	}
}

var (
	EndpointsFromDB StoredEndpoints
	//EndpointsForReq FastestEndpoints
)

func StoreEndpointsFromDB() {
	log.Println("Started storing!")
	for {
		entries, err := db.GetAll()
		if err != nil {
			log.Println(err)
			continue
		}

		for _, j := range entries {
			EndpointsFromDB.Add(j)
		}

		time.Sleep(time.Minute * 10)
	}
}

//func SetFastestEndpoints(){
//	log.Println("Started set fastest endpoints!")
//	for {
//		time.Sleep(time.Minute * 10)
//	}
//}

func GetEndpoint(currency string) (string, error) {
	endpoints := EndpointsFromDB.GetByCurrency(currency)
	if endpoints == nil {
		return "", errors.New("Not found")
	}

	addresses := endpoints.Addresses

	rand.Seed(time.Now().UnixNano())

	result := addresses[rand.Intn(len(addresses))]

	return result, nil
}
