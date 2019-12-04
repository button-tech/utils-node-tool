package handlers

import (
	"github.com/button-tech/logger"
	"github.com/button-tech/utils-node-tool/nodetools"
	b "github.com/button-tech/utils-node-tool/nodetools"
	"github.com/button-tech/utils-node-tool/types/responses"
	"github.com/qiangxue/fasthttp-routing"
	"os"
	"time"
)

func GetBalance(c *routing.Context) error {

	start := time.Now()

	address := c.Param("address")

	balance, err := b.GetUtxoBasedBalance(address)
	if err != nil {
		logger.HandlerError("GetBalance", err)
		return err
	}

	response := responses.BalanceResponse{Balance: balance}

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	logger.LogRequest(time.Since(start), os.Getenv("BLOCKCHAIN"), "GetBalance", false)

	return nil
}

func GetUtxo(c *routing.Context) error {

	start := time.Now()

	address := c.Param("address")

	utxoArray, err := nodetools.GetUtxo(address)
	if err != nil {
		logger.HandlerError("GetUtxo", err)
		return err
	}

	response := responses.UTXOResponse{Utxo: utxoArray}

	if err := responses.JsonResponse(c, response); err != nil {
		return err
	}

	logger.LogRequest(time.Since(start), os.Getenv("BLOCKCHAIN"), "GetUtxo", false)

	return nil
}
