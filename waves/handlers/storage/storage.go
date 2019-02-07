package storage

import "os"

var WavesURL = os.Getenv("WAVES_NODE")

type BalanceData struct {
	Balance int64 `json:"balance"`
}


