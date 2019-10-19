package schema

type EndpointsData struct {
	Currency  string   `json:"currency"`
	Addresses []string `json:"addresses"`
	Port      int      `json:"port"`
	Stopped   []string `json:"stopped"`
	Reserve   string   `json:"reserve"`
}
