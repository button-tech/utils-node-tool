package requests

type BalancesRequest struct {
	Addresses []string `json:"addresses"`
}

type UtxoBasedBalance struct {
	Balance string `json:"balance"`
}

type UtxoBasedBlocksHeight struct {
	Backend struct {
		Blocks int64 `json:"blocks"`
	}
}

type EthBasedBlocksHeight struct {
	Result struct {
		Number string `json:"number"`
	}
}

type EthEstimateGasRequest struct {
	ContractAddress string `json:"contractAddress"`
	Data            string `json:"data"`
}
