package configs

// Validator is a struct that represents what CAF should
// check for in a reqeust
type Validator struct {
	Monitor string `json:"monitor"`
}
