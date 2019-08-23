package proxy

import (
	"log"

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
		CertFile: config.Cert,
		KeyFile:  config.Key,
	}, func(err error) {
		log.Print(err.Error())
	})

	// Print all routes that were hit
	mm.All("/.*", middleman.RouteLogger())

	// Read request body and store it in store.Body
	mm.All("/.*", middleman.BodyReader())

	// ======================== Proxy code begins here ========================

	//mm.All("/.*", PrintRequestBody())

	// ========================= Proxy code ends here =========================

	// Forward request to the target
	mm.All("/.*", SendRequest(config.Target))

	// Print the target response body
	//mm.All("/.*", PrintTargetResponseBody())

	// Forward response to the client
	mm.All("/.*", SendResponse())

	log.Println("[Middleman is listening on]:", config.Addr)
	log.Println("[Proxy is fowrarding to]:", config.Target)

	// Begin listening
	err := mm.ListenAndServeTLS()

	// If an error occured, print a message
	if err != nil {
		log.Fatalln("[Failed creating a server]:", err)
	}
}
