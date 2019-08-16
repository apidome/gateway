package proxy

import (
	"log"

	"github.com/Creespye/caf/internal/pkg/middleman"
	mmMiddlewares "github.com/Creespye/caf/internal/pkg/middleman/middlewares"
	proxyMiddlewares "github.com/Creespye/caf/internal/pkg/proxy/middlewares"
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

	// Print all routes that were hit
	mm.Use(mmMiddlewares.RouteLogger())

	// Read request body and store it in store.Body
	mm.Use(mmMiddlewares.BodyReader())

	// Forward request to the target
	mm.Use(proxyMiddlewares.ForwardRequest(config.Target))

	log.Println("[Middleman is listening on]:", config.Addr)

	// Begin listening
	err := mm.ListenAndServeTLS()

	// If an error occured, print a message
	if err != nil {
		log.Fatalln("[Failed creating a server]:", err)
	}
}
