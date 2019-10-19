package nodetools

import (
	"github.com/button-tech/utils-node-tool/db"
	"log"
)

func Max(array []int64) int64 {
	var max int64 = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
	}
	return max
}

func DeleteEntry(currency, address string) error {
	err := db.AddToStoppedList(currency, address)
	if err != nil {
		return err
	}

	log.Printf("Add to stopped list %s: %s", currency, address)

	return nil
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
