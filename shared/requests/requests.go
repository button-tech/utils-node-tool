package requests

type BalancesRequest struct {
	Addresses []string `json:"addresses"`
}

type EthEstimateGasRequest struct {
	ToAddress    string `json:"toAddress"`
	TokenAddress string `json:"tokenAddress"`
	Amount       string `json:"amount"`
}
