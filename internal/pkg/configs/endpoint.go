package configs

type Endpoint struct {
	Path   string `json:"path"`
	Method string `json:"method"`
	Schema string `json:"schema"`
}
