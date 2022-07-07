package entities

type UsageStatus struct {
	APIKey        string   `json:"-"`
	Status        string   `json:"status"`
	CreationTime  randTime `json:"creationTime"`
	TotalRequests uint64   `json:"totalRequests"`
	TotalBits     uint64   `json:"totalBits"`
	RequestsLeft  uint64   `json:"requestsLeft"`
	BitsLeft      uint64   `json:"bitsLeft"`
}
