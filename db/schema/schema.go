package schema

type EndpointsData struct {
	Currency  string   `json:"currency"`
	Addresses []string `json:"addresses"`
	Stopped   []string `json:"stopped"`
}
