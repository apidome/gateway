package configs

/*
A struct that hold the configuration of the trusted side.
*/
type In struct {
	Targets []Target `json:"targets"`
}
