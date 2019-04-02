package responses

type UTXO struct {
	Address       string  `json:"address"`
	Txid          string  `json:"txid"`
	Vout          int     `json:"vout"`
	ScriptPubKey  string  `json:"scriptPubKey"`
	Amount        float64 `json:"amount"`
	Satoshis      int     `json:"satoshis"`
	Height        int     `json:"height"`
	Confirmations int     `json:"confirmations"`
}

type BalanceResponse struct {
	Balance string `json:"balance" example:"0"`
}

type TransactionFeeResponse struct {
	Fee float64 `json:"fee" example:"0"`
}

type UTXOResponse struct {
	Utxo []UTXO `json:"utxo"`
}

type BalancesResponse struct {
	Balances map[string]string `json:"balances"`
}
