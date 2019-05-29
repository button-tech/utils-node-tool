package schema

type Endpoints struct {
	Currency  string   `json:"currency"`
	Addresses []string `json:"addresses"`
	Port      int      `port`
	Reserve   string   `json:"reserve"`
}
