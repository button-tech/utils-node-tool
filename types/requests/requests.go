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

type XRPTxData struct {
	TxBlob string `json:"txBlob"`
	Devnet bool   `json:"devnet"`
}

type XRPSentTxInfo struct {
	Result struct {
		EngineResult        string `json:"engine_result"`
		EngineResultCode    int    `json:"engine_result_code"`
		EngineResultMessage string `json:"engine_result_message"`
		Status              string `json:"status"`
		TxJSON              struct {
			Fee  string `json:"Fee"`
			Hash string `json:"hash"`
		} `json:"tx_json"`
	} `json:"result"`
}

type XRPTxToSubmit struct {
	Method string `json:"method"`
	Params []struct {
		TxBlob string `json:"tx_blob"`
	} `json:"params"`
}
