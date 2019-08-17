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

	// Print all routes that were hit
	mm.All("/.*", middleman.RouteLogger())

	// Read request body and store it in store.Body
	mm.All("/.*", middleman.BodyReader())

	// ======================== Proxy code begins here ========================

	mm.Get("/pp", func(res http.ResponseWriter, req *http.Request,
		store *middleman.Store, end middleman.End) {

		log.Println(store.RequestBody)
	})

	mm.Get("/pp", func(res http.ResponseWriter, req *http.Request,
		store *middleman.Store, end middleman.End) {

		log.Println(store.TargetResponseBody)
	})

	mm.All("/.*", PrintRequestBody())

	mm.All("/.*", PrintTargetResponseBody())

	// ========================= Proxy code ends here =========================

	// Forward request to the target
	mm.All("/.*", SendRequest(config.Target))

	// Forward response to the client
	mm.All("/.*", SendResponse())

	log.Println("[Middleman is listening on]:", config.Addr)

	// Begin listening
	err := mm.ListenAndServeTLS()

	// If an error occured, print a message
	if err != nil {
		log.Fatalln("[Failed creating a server]:", err)
	}
}
