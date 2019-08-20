package json_configs

/*
A struct that hold the configuration of the untrusted side.
*/
type Out struct {
	Port string `json:"port"`
	SSL  bool   `json:"ssl"`
}
