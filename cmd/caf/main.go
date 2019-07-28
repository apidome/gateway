package main

import "../../internal/pkg/proxy"

func main() {
	proxy.Start(proxy.Config{
		Addr:   ":8080",
		Target: "google.com",
		Cert:   "../../configs/certs/localhost/localhost.cert",
		Key:    "../../configs/certs/localhost/localhost.key",
	})
}
