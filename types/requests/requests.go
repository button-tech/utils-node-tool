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

type BnBBalance struct {
	AccountNumber int    `json:"account_number"`
	Address       string `json:"address"`
	Balances      []struct {
		Free   string `json:"free"`
		Frozen string `json:"frozen"`
		Locked string `json:"locked"`
		Symbol string `json:"symbol"`
	} `json:"balances"`
	Sequence int `json:"sequence"`
}

type CosmosBalance []struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type AlgorandBalance struct {
	Participation struct {
		Vrfpkb64  string `json:"vrfpkb64"`
		Partpkb64 string `json:"partpkb64"`
		Votefst   int    `json:"votefst"`
		Votelst   int    `json:"votelst"`
		Votekd    int    `json:"votekd"`
	} `json:"participation"`
	Amount                      uint64 `json:"amount"`
	Pendingrewards              int    `json:"pendingrewards"`
	Address                     string `json:"address"`
	Assets                      string `json:"assets"`
	Round                       int    `json:"round"`
	Thisassettotal              string `json:"thisassettotal"`
	Amountwithoutpendingrewards int    `json:"amountwithoutpendingrewards"`
	Rewards                     int    `json:"rewards"`
	Status                      string `json:"status"`
}
