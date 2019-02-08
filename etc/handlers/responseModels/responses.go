package responses

type BalanceResponse struct {
	Balance string `json:"balance"`
}

type TransactionFeeResponse struct {
	Fee float64 `json:"fee" example:"0"`
}

type GasPriceResponse struct {
	GasPrice int64 `json:"gasPrice" example:"0"`
}

type BalancesResponse struct {
	Balances map[string]string `json:"balances"`
}
