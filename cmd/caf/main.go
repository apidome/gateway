package main

import "github.com/Creespye/caf/internal/pkg/proxy"

func main() {
	proxy.Start(proxy.Config{
		Addr:   "localhost:8080",
		Target: "https://google.com",
		Cert:   "../../configs/certs/localhost/localhost.cert",
		Key:    "../../configs/certs/localhost/localhost.key",
	})
}
