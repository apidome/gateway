package middleman

import (
	"net/http"
	"strings"
)

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

// route is a struct that holds all middlewares for a given route
type route struct {
	path        string
	middlewares map[string][]Middleware
}

// Middleware is the function needed to implement as a middleware
type Middleware func(res http.ResponseWriter, req *http.Request, store map[string]string, end End)

// End is the function that will be called to break the continuation of middlewares
type End func()

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

// addMiddleware Adds a middleware of a certain method to a route
func (mm *Middleman) addMiddleware(path string, method string, handler Middleware) {
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
		newMiddlewares := make(map[string][]Middleware)

		newMiddlewares[method] = []Middleware{
			handler,
		}

		mm.routes = append(mm.routes, route{
			path:        path,
			middlewares: newMiddlewares,
		})
	}
}

// Get Adds a GET middleware to a route
func (mm *Middleman) Get(path string, handler Middleware) {
	mm.addMiddleware(path, "GET", handler)
}

// Post Adds a POST middleware to a route
func (mm *Middleman) Post(path string, handler Middleware) {
	mm.addMiddleware(path, "POST", handler)
}

// mainHandler is the main function that receives all requests and calls the
// correct middlewares
func (mm *Middleman) mainHandler(res http.ResponseWriter, req *http.Request) {
	store := make(map[string]string)

	// Find all routes on the way to the desired route
	paths := strings.Split(req.RequestURI, "/")

	// Remove the empty string at the end of the paths array
	// (Split returns it and its useless)
	if paths[len(paths)-1] == "" {
		paths = paths[0:1]
	}

	// Prefix all sub-paths with a '/'
	for i := range paths {
		paths[i] = "/" + paths[i]
	}

	// Define current path as an empty string and concatenate sub paths to it
	// during iteration
	currentPath := ""

	// Iterate over all sub paths of this request
	for _, path := range paths {
		// Concat the paths together
		currentPath += path

		// Remove '//' because each path is prefixed with a '/'
		currentPath = strings.ReplaceAll(currentPath, "//", "/")

		// Iterate over all routes and execute middleware of all sub paths
		for _, route := range mm.routes {
			if route.path == currentPath {
				for _, middleware := range route.middlewares[req.Method] {
					cont := true

					middleware(res, req, store, func() {
						cont = false
					})

					// If the end function was called, break middleware execution
					if !cont {
						break
					}
				}
			}
		}
	}
}
