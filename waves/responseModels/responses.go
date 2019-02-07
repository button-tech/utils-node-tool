package responses

type BalanceResponse struct {
	Balance string `json:"balance"`
}

type BalancesResponse struct {
	Balances map[string]string `json:"balances"`
}
