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
	EthEdnpoints,
	EtcEndpoints schema.Endpoints
}

var EndpointsFromDB StoredEndpoints

func (s *StoredEndpoints) Add(entry schema.Endpoints) {
	s.Lock()
	switch entry.Currency {
	case "btc":
		s.BtcEndpoints = entry
	case "ltc":
		s.LtcEndpoints = entry
	case "eth":
		s.EthEdnpoints = entry
	case "etc":
		s.EtcEndpoints = entry
	}
	s.Unlock()
}

func (s *StoredEndpoints) Get(currency string) *schema.Endpoints {
	s.RLock()
	defer s.RUnlock()
	switch currency {
	case "btc":
		return &s.BtcEndpoints
	case "ltc":
		return &s.LtcEndpoints
	case "eth":
		return &s.EthEdnpoints
	case "etc":
		return &s.EtcEndpoints
	}
	return nil
}

func StoreEndpoints() {
	log.Println("Start storing!")
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

func GetEndpoint(currency string) (string, error) {
	endpoints := EndpointsFromDB.Get(currency)
	if endpoints == nil {
		return "", errors.New("Not found")
	}

	addresses := endpoints.Addresses

	rand.Seed(time.Now().UnixNano())

	result := addresses[rand.Intn(len(addresses))]

	return result, nil
}
