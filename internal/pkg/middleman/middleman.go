package middleman

import (
	"net/http"
)

// route is a struct that holds all middlewares for a given route
type route struct {
	path        string
	middlewares map[string][]http.HandlerFunc
}

// Middleman is a struct that holds all middlewares
type Middleman struct {
	config Config
	routes []route
}

// Config holds the configurations for the underlying web server
type Config struct {
	Addr     string
	Target   string
	CertFile string
	KeyFile  string
}

// NewMiddleman returns a new instance of a middleman
func NewMiddleman(config Config) Middleman {
	return Middleman{
		config: config,
	}
}

// ListenAndServeTLS starts the https server
func (mm *Middleman) ListenAndServeTLS() error {
	http.HandleFunc("/", mm.mainHandler)

	// Start the listener, and if an error occures, pass is up to the caller
	err := http.ListenAndServeTLS(mm.config.Addr, mm.config.CertFile, mm.config.KeyFile, nil)

	return err
}

func (mm *Middleman) addMiddleware(path string, method string, handler http.HandlerFunc) {
	foundRoute := false

	// Tries to find the route in the middleman struct to add a new middleware to it
	for _, route := range mm.routes {
		if route.path == path {
			foundRoute = true

			route.middlewares[method] = append(route.middlewares[method], handler)
		}
	}

	// If the route was not found, create it in the middleman
	if !foundRoute {
		newMiddlewares := make(map[string][]http.HandlerFunc)

		newMiddlewares[method] = []http.HandlerFunc{
			handler,
		}

		mm.routes = append(mm.routes, route{
			path:        path,
			middlewares: newMiddlewares,
		})
	}
}

// Get Adds a GET middleware to a route
func (mm *Middleman) Get(path string, handler http.HandlerFunc) {
	mm.addMiddleware(path, "GET", handler)
}

// Post Adds a POST middleware to a route
func (mm *Middleman) Post(path string, handler http.HandlerFunc) {
	mm.addMiddleware(path, "POST", handler)
}

// mainHandler is the main function that receives all requests and calls the
// correct middlewares
func (mm *Middleman) mainHandler(res http.ResponseWriter, req *http.Request) {
	for _, route := range mm.routes {
		if route.path == req.RequestURI {
			for _, middleware := range route.middlewares[req.Method] {
				middleware(res, req)
			}
		}
	}
}
