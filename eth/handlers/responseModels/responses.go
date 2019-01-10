package responses

type BalanceResponse struct {
	Balance float64 `json:"balance" example:"0"`
}

type TransactionFeeResponse struct {
	Fee float64 `json:"fee" example:"0"`
}

type GasPriceResponse struct {
	GasPrice int64 `json:"gasPrice" example:"0"`
}

type TokenBalanceResponse struct {
	TokenBalance float64 `json:"tokenBalance" example:"0"`
}

type BalancesResponse struct {
	Balances map[string]float64 `json:"balances"`
}
