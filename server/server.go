package main

import (
	"../handlers/eth"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	// get ETH/ETC acc balance
	r.GET("/balance/:eth/:address", eth.GetBalance)

	// for ETH/ETC fee
	r.GET("/fee/:eth", eth.GetTxFee)

	// ETH/ETC gas price
	r.GET("/gasprice/:eth", eth.GetGasPrice)

	// for ETH only
	// get ETH acc token balance
	// example - /tbalance/knc/*
	r.GET("/tbalance/:token/:address", eth.GetTokenBalance)

	// send rawTx ETH/ETC
	// example {"raw_tx":"f86d8202b284773594008252..."}
	r.POST("/send_tx/:eth", eth.SendTX)

	r.Run()
}
