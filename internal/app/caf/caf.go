package caf

import "github.com/Creespye/caf/internal/pkg/proxy"

// Start starts CAF
func Start() {
	proxy.Start(proxy.Config{
		Addr:   "localhost:8080",
		Target: "https://google.com",
		Cert:   "../../configs/certs/localhost/localhost.cert",
		Key:    "../../configs/certs/localhost/localhost.key",
	})
}
