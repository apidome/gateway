package proxy

import (
	"log"
	"net/http"

	"github.com/Creespye/caf/internal/pkg/middleman"
)

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
	// Creating a new middleman (middleware manager)
	mm := middleman.NewMiddleman(middleman.Config{
		Addr:     config.Addr,
		Target:   config.Target,
		CertFile: config.Cert,
		KeyFile:  config.Key,
	})

	msg := ""

	// Create a new middleware for the root route
	mm.Get("/", func(res http.ResponseWriter, req *http.Request) {
		msg = ""
		msg += "Im first!\n"
	})

	// Add another middleware to the root route
	mm.Get("/", func(res http.ResponseWriter, req *http.Request) {
		msg += "Im second!\n"

		res.Write(([]byte)(msg))
	})

	// Begin listening
	err := mm.ListenAndServeTLS()

	// If an error occured, print a message
	if err != nil {
		log.Fatalln("Failed creating a server: ", err)
	} else {
		log.Println("Middleman listening on ", config.Addr)
	}
}
