package configs

const (
	TypeRest = "REST"
)

// Declaration of api object
type Api struct {
	Type      string      `json:"type"`
	Version   string      `json:"version"`
	Validator Validator   `json:"validator"`
	Endpoints []*Endpoint `json:"endpoints"`
}
