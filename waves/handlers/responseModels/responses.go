package responses

type BalanceResponse struct {
	Balance string `json:"balance"`
}

type BalancesResponse struct {
	Balances map[string]string `json:"balances"`
}

type TransactionFeeResponse struct {
	Fee float64 `json:"fee" example:"0"`
}
