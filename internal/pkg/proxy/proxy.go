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
		CertFile: config.Cert,
		KeyFile:  config.Key,
	}, func(err error) bool {
		log.Print("[Middleware Error]: " + err.Error())

		return false
	})

	// Print all routes that were hit
	mm.All("/.*", middleman.RouteLogger())

	// Read request body and store it in store.Body
	mm.All("/.*", middleman.BodyReader())

	// ==================== Request proxy code begins here ====================

	// ===================== Request proxy code ends here =====================

	// Create the target request
	mm.All("/.*", CreateTargetRequest(config.Target))

	// Forward request to the target
	mm.All("/.*", SendTargetRequest())

	// Read the target response body from store.TargetResponse
	mm.All("/.*", ReadTargetResponseBody())

	// ==================== Response proxy code begins here ===================

	// Change referer header for youtube
	mm.All("/.*", func(res http.ResponseWriter, req *http.Request,
		store *middleman.Store, end middleman.End) error {

		//store.TargetRequest.Header.Del("Referer")

		return nil
	})

	// ===================== Response proxy code ends here ====================

	// Forward response to the client
	mm.All("/.*", SendTargetResponse())

	log.Println("[Middleman is listening on]:", config.Addr)
	log.Println("[Proxy is fowrarding to]:", config.Target)

	// Begin listening
	err := mm.ListenAndServeTLS()

	// If an error occured, print a message
	if err != nil {
		log.Fatalln("[Failed creating a server]:", err)
	}
}
