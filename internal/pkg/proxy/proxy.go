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

// Start starts the proxy server and begins operating on requests
func Start(config Config) {
	// Creating a new middleman (middleware manager)
	mm := middleman.NewMiddleman(middleman.Config{
		Addr:     config.Addr,
		Target:   config.Target,
		CertFile: config.Cert,
		KeyFile:  config.Key,
	})

	// Create a new middleware for the root route
	mm.Get("/", func(res http.ResponseWriter, req *http.Request,
		store map[string]string, end middleman.End) {
	})

	// Add another middleware to the root route
	mm.Get("/hey", func(res http.ResponseWriter, req *http.Request,
		store map[string]string, end middleman.End) {
	})

	mm.Get("/hey/im", func(res http.ResponseWriter, req *http.Request,
		store map[string]string, end middleman.End) {
	})

	mm.Get("/hey/im/omer", func(res http.ResponseWriter, req *http.Request,
		store map[string]string, end middleman.End) {
	})

	mm.Use(func(res http.ResponseWriter, req *http.Request,
		store map[string]string, end middleman.End) {
	})

	mm.All("/", func(res http.ResponseWriter, req *http.Request,
		store map[string]string, end middleman.End) {
	})

	// Begin listening
	err := mm.ListenAndServeTLS(func() {
		log.Println("Middleman is listening on", config.Addr)
	})

	// If an error occured, print a message
	if err != nil {
		log.Fatalln("Failed creating a server: ", err)
	}
}
