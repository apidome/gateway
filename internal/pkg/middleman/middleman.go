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

// Store is a struct that holds data between middlewares
type Store struct {
	RequestBody        []byte
	TargetResponse     *http.Response
	TargetResponseBody []byte
	Generics           map[string]interface{}
}

// Middleware is the function needed to implement as a middleware
type Middleware func(res http.ResponseWriter, req *http.Request,
	store *Store, end End)

// End is the function that will be called to break the continuation of middlewares
type End func()

var (
	methods = []string{
		http.MethodConnect,
		http.MethodDelete,
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
		http.MethodTrace,
	}
)

// NewMiddleman returns a new instance of a middleman
func NewMiddleman(config Config) Middleman {
	return Middleman{
		config: config,
	}
}

// ListenAndServeTLS starts the https server
func (mm *Middleman) ListenAndServeTLS() error {
	http.HandleFunc("/", mm.mainHandler)

	// Start the listener, and if an error occures, pass it up to the caller
	err := http.ListenAndServeTLS(mm.config.Addr, mm.config.CertFile, mm.config.KeyFile, nil)

	return err
}

// mainHandler is the main function that receives all requests and calls the correct middlewares
func (mm *Middleman) mainHandler(res http.ResponseWriter, req *http.Request) {
	// Create a store to hold information between middlewares
	store := Store{}

	// TODO Rework this if rquired
	// Execute generic middlewares ('Use' middlewares)
	//mm.runMiddlewares("/", req.Method, res, req, store)

	// Find all paths on the way to the desired path
	paths := strings.Split(req.RequestURI, "/")

	// Remove the empty string at the end of the paths array
	// (Split returns it and its useless)
	if paths[len(paths)-1] == "" {
		paths = paths[0 : len(paths)-1]
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

		// Execute middlewares of the current route
		cont := mm.runMiddlewares(currentPath, req.Method, res, req, &store)

		if !cont {
			break
		}
	}
}

// addMiddleware Adds a middleware of a certain method to a route
func (mm *Middleman) addMiddleware(path string, method string, middleware Middleware) {
	foundRoute := false

	// Tries to find the route in the middleman struct to add a new middleware to it
	for _, route := range mm.routes {
		if route.path == path {
			foundRoute = true

			route.middlewares[method] = append(route.middlewares[method], middleware)
		}
	}

	// If the route was not found, create it in the middleman
	if !foundRoute {
		newMiddlewares := make(map[string][]Middleware)

		newMiddlewares[method] = []Middleware{
			middleware,
		}

		mm.routes = append(mm.routes, route{
			path:        path,
			middlewares: newMiddlewares,
		})
	}
}

// runMiddlewares executes add middlewares of a specific path (and all of its sub-paths)
// it returns a boolean indication if middleware execution should continue
func (mm *Middleman) runMiddlewares(path string, method string,
	res http.ResponseWriter,
	req *http.Request,
	store *Store) bool {

	// Declare a variable to indicate if execution should be terminated before all
	// middlewares were executed
	terminate := false

	// Iterate over all routes and execute middleware of all sub paths
	for _, route := range mm.routes {
		if route.path == path {
			for _, middleware := range route.middlewares[method] {
				middleware(res, req, store, func() {
					terminate = true
				})

				// If the end function was called, break middleware execution
				if terminate {
					break
				}
			}

			if terminate {
				break
			}
		}
	}

	return !terminate
}

// Get Adds a GET middleware to a route
func (mm *Middleman) Get(path string, middleware Middleware) {
	mm.addMiddleware(path, http.MethodGet, middleware)
}

// Head Adds a HEAD middleware to a route
func (mm *Middleman) Head(path string, middleware Middleware) {
	mm.addMiddleware(path, http.MethodHead, middleware)
}

// Post Adds a POST middleware to a route
func (mm *Middleman) Post(path string, middleware Middleware) {
	mm.addMiddleware(path, http.MethodPost, middleware)
}

// Put Adds a PUT middleware to a route
func (mm *Middleman) Put(path string, middleware Middleware) {
	mm.addMiddleware(path, http.MethodPut, middleware)
}

// Delete Adds a DELETE middleware to a route
func (mm *Middleman) Delete(path string, middleware Middleware) {
	mm.addMiddleware(path, http.MethodDelete, middleware)
}

// Connect Adds a CONNECT middleware to a route
func (mm *Middleman) Connect(path string, middleware Middleware) {
	mm.addMiddleware(path, http.MethodConnect, middleware)
}

// Options Adds a OPTIONS middleware to a route
func (mm *Middleman) Options(path string, middleware Middleware) {
	mm.addMiddleware(path, http.MethodOptions, middleware)
}

// Trace Adds a TRACE middleware to a route
func (mm *Middleman) Trace(path string, middleware Middleware) {
	mm.addMiddleware(path, http.MethodTrace, middleware)
}

// Patch Adds a PATCH middleware to a route
func (mm *Middleman) Patch(path string, middleware Middleware) {
	mm.addMiddleware(path, http.MethodPatch, middleware)
}

// All Adds a middleware to all methods of a route
func (mm *Middleman) All(path string, middleware Middleware) {
	for _, method := range methods {
		mm.addMiddleware(path, method, middleware)
	}
}

// Use Adds a generic middleware to the root path of the listener (and any sub-paths)
func (mm *Middleman) Use(middleware Middleware) {
	mm.All("/", middleware)
}
