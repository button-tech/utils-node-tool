package responses


type BalanceResponse struct {
	Balance    float64 `json:"balance" example:"0"`
}

type TransactionFeeResponse struct {
	Fee float64 `json:"fee" example:"0"`
}