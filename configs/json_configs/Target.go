package json_configs

/*
This struct represents a Target json object.
*/
type Target struct {
	Host       string      `json:"host"`
	Port       string      `json:"port"`
	SSL        bool        `json:"ssl"`
	ClientAuth bool        `json:"clientAuth"`
	Validators []Validator `json:"validators"`
}

/*
A function that returns a string of the URL of the target.
*/
func (t Target) GetURL() string {
	var scheme string

	// Check if the target is listening on http or https.
	if t.SSL {
		scheme = "https://"
	} else {
		scheme = "http://"
	}

	// Return the full URL.
	return scheme + t.Host + ":" + t.Port
}
