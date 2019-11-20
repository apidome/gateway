package configs

const (
	// TypeRest Indicates REST configurations
	TypeRest = "REST"
)

// API holds information on a specific API
type API struct {
	Type      string      `json:"type"`
	Version   string      `json:"version"`
	Validator Validator   `json:"validator"`
	Endpoints []*Endpoint `json:"endpoints"`
}
