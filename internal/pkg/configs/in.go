package configs

//In is a struct that hold the configuration of the trusted side.
type In struct {
	Targets []Target `json:"targets"`
}
