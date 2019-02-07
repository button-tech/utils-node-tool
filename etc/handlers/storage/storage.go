package storage

import (
	"github.com/onrik/ethrpc"
	"os"
)

var (
	EtcURL = os.Getenv("ETC_NODE")

	EtcClient = ethrpc.New(EtcURL)
)
