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

type UtxoBasedTxOutputs struct {
	Vout []struct {
		ScriptPubKey struct {
			Hex       string   `json:"hex"`
			Addresses []string `json:"addresses"`
		}
	} `json:"vout"`
}

type EthEstimateGasRequest struct {
	ContractAddress string `json:"contractAddress"`
	Data            string `json:"data"`
}

type TezosBalance []string

type XRPBalance struct {
	Result      string `json:"result"`
	LedgerIndex int    `json:"ledger_index"`
	Limit       int    `json:"limit"`
	Balances    []struct {
		Currency string `json:"currency"`
		Value    string `json:"value"`
	} `json:"balances"`
}
