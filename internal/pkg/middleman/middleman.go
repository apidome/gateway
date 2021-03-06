package middleman

import (
	"crypto/tls"
	"net/http"
	"regexp"
)

// Store is a struct that holds data between middlewares
type Store map[string]interface{}

// Middleware is the function needed to implement as a middleware
type Middleware func(res http.ResponseWriter, req *http.Request,
	store Store, end End) error

type errorHandler func(path, method string, err error) bool

// handler is a struct that hold middleware information
type middlewareHandler struct {
	middleware Middleware
	path       string
	method     string
}

// Middleman is a struct that holds all middlewares
type Middleman struct {
	handlers     []middlewareHandler
	errorHandler errorHandler
	httpServer   http.Server
}

// End is the function that will be called to break
// the continuation of middlewares
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

// InitMiddleman initializes a middleman instance
func InitMiddleman(mm *Middleman, addr string, errHandler errorHandler) {
	// Disable HTTP/2
	tlsNextProto := make(map[string]func(*http.Server, *tls.Conn, http.Handler))

	mm.errorHandler = errHandler

	mm.httpServer.Addr = addr
	mm.httpServer.TLSNextProto = tlsNextProto
	mm.httpServer.Handler = http.HandlerFunc(mm.mainHandler)
}

// NewMiddleman returns a new instance of a middleman
func NewMiddleman(addr string, errHandler errorHandler) *Middleman {
	mm := &Middleman{}

	InitMiddleman(mm, addr, errHandler)

	return mm
}

// ListenAndServeTLS starts the https server
func (mm *Middleman) ListenAndServeTLS(certFile, keyFile string) error {
	err := mm.httpServer.ListenAndServeTLS(certFile, keyFile)

	return err
}

// ListenAndServe starts the http server
func (mm *Middleman) ListenAndServe() error {
	err := mm.httpServer.ListenAndServe()

	return err
}

// emitError calls the error handler callback to inform the user of an error
// and returns if execution should continue
func (mm *Middleman) emitError(path, method string, err error) bool {
	if mm.errorHandler != nil {
		return mm.errorHandler(path, method, err)
	}

	// If no error handler was configured, do not stop execution
	return true
}

// mainHandler is the main function that receives all
// requests and calls the correct middlewares
func (mm *Middleman) mainHandler(res http.ResponseWriter, req *http.Request) {
	// Store holds data between middlewares
	store := Store{}

	_, err := mm.runMiddlewares(res, req, store)

	if err != nil {
		mm.emitError(req.URL.Path, req.Method, err)
	}
}

// addMiddleware adds a middleware to the middleware store
func (mm *Middleman) addMiddleware(path string, method string,
	middleware Middleware) error {
	// We are using the path argument as a regular expression, so in order
	// to fit our needs we surround it with ^ and $ to avoid regex matching
	// anything that contains this path, rather than beginning with it or
	// being equal to it
	regexPath := "^" + path + "$"

	_, err := regexp.Compile(regexPath)

	if err != nil {
		return err
	}

	mm.handlers = append(mm.handlers, middlewareHandler{
		middleware,
		regexPath,
		method,
	})

	return nil
}

// runMiddlewares runs middlewares on a request
// Returns a bool value to indicate if execution stopped
// Returns an error if any occured
func (mm *Middleman) runMiddlewares(res http.ResponseWriter, req *http.Request,
	store Store) (bool, error) {
	// Indication wether execution should be stopped
	cont := true

	// Define the end function
	end := func() {
		cont = false
	}

	for _, handler := range mm.handlers {
		// If the middleware called the end function, middleware execution
		// should be stopped
		if !cont {
			break
		}

		// Match the regex of the handler to the request's uri path
		regexMatch, err := regexp.MatchString(handler.path,
			req.URL.Path)

		if err != nil {
			return false, err
		}

		if regexMatch && handler.method == req.Method {
			err := handler.middleware(res, req, store, end)

			// If an error occured in the middleware, emit the error
			if err != nil {
				// Raise error emitter and decide to continue or break
				continueAfterError :=
					mm.emitError(req.URL.Path, req.Method, err)

				if !continueAfterError {
					break
				}
			}
		}
	}

	return cont, nil
}
