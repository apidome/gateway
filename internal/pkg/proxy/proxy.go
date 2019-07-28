package proxy

import "net/http"
import "log"

// Config is a struct that holds all configurations of the proxy server
type Config struct {
	Addr   string
	Target string
	Cert   string
	Key    string
}

var proxyConf Config

// Start starts the proxy server and begins operating on requests
func Start(config Config) {
	proxyConf = config
	http.HandleFunc("/", forwardHandler)
	err := http.ListenAndServeTLS(config.Addr, config.Cert, config.Key, nil)

	if err != nil {
		log.Fatal("ListenAndServerTLS: ", err)
	} else {
		log.Println("Proxy is up")
	}
}

func forwardHandler(res http.ResponseWriter, req *http.Request) {

	// Read Request

	// Create target request
	switch req.Method {
	case "GET":
		targetRes, err := http.Get(proxyConf.Target)

		if err != nil {
			log.Fatal("Request to target failed:", err)
		} else {
			targetRes.Write(res)
		}
	}
	// Forward target response as response
	res.Write([]byte(req.Method))
}
