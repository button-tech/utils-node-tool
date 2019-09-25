package responses

import (
	"encoding/json"
	"github.com/qiangxue/fasthttp-routing"
)

type UTXO struct {
	Address       string `json:"address,omitempty"`
	Txid          string `json:"txid"`
	Vout          int    `json:"vout"`
	ScriptPubKey  string `json:"scriptPubKey,omitempty"`
	Amount        string `json:"amount"`
	Satoshis      int    `json:"satoshis"`
	Height        int    `json:"height"`
	Confirmations int    `json:"confirmations"`
	LegacyAddress string `json:"legacyAddress,omitempty"`
	CashAddress   string `json:"cashAddress,omitempty"`
}

type UTXOResponse struct {
	Utxo []UTXO `json:"utxo"`
}

type BalanceResponse struct {
	Balance string `json:"balance" example:"0"`
}

type TransactionFeeResponse struct {
	Fee float64 `json:"fee" example:"0"`
}

type BalancesResponse struct {
	Balances map[string]string `json:"balances"`
}

type GasPriceResponse struct {
	GasPrice int64 `json:"gasPrice" example:"0"`
}

type BalanceData struct {
	Balance int64 `json:"balance"`
}

type GasLimitResponse struct {
	GasLimit uint64 `json:"gasLimit"`
}

func JsonResponse(ctx *routing.Context, data interface{}) error {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, HEAD")
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	ctx.Response.Header.SetCanonical([]byte("Content-Type"), []byte("application/json"))
	ctx.Response.SetStatusCode(200)
	if err := json.NewEncoder(ctx).Encode(data); err != nil {
		return err
	}
	return nil
}
