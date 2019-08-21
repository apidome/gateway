package configs

/*
A struct that hold the configuration of the untrusted side.
*/
type Out struct {
	Port            string `json:"port"`
	SSL             bool   `json:"ssl"`
	CertificatePath string `json:"certPath"`
	KeyPath         string `json:"keyPath"`
}
