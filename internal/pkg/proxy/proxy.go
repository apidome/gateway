package proxy

import "net/http"

// ProxyConfig is a struct that holds all configurations of the proxy server
type ProxyConfig struct {
	addr   string
	target string
	cert   string
	key    string
}

var conf ProxyConfig

func Start(addr string,
	target string,
	cert string,
	key string,
	handler func(res http.ResponseWriter, req *http.Request)) {

	conf = ProxyConfig{
		addr:   addr,
		target: target,
		cert:   cert,
		key:    key,
	}

	http.HandleFunc("/", forwardHandler)
}

func forwardHandler(res http.ResponseWriter, req *http.Request) {

	// Read Request

	// Create target request

	// Forward target response as response
}
