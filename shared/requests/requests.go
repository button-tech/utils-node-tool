package requests

type BalancesRequest struct {
	Addresses []string `json:"addresses"`
}

type EthEstimateGasRequest struct {
	ContractAddress string `json:"contractAddress"`
	Data            string `json:"data"`
}
