package middleman

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"regexp"
)

// Store is a struct that holds data between middlewares
type Store map[string]interface{}

// Middleware is the function needed to implement as a middleware
type Middleware func(res http.ResponseWriter, req *http.Request,
	store Store, end End) error

// handler is a struct that hold middleware information
type middlewareHandler struct {
	middleware Middleware
	path       string
	method     string
}

// Middleman is a struct that holds all middlewares
type Middleman struct {
	handlers     []middlewareHandler
	errorHandler func(error) bool
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

// NewMiddleman returns a new instance of a middleman
func NewMiddleman(mm *Middleman,
	addr string,
	errorHandler func(error) bool) {

	// Disable HTTP/2
	tlsNextProto := make(map[string]func(*http.Server, *tls.Conn, http.Handler))

	*mm = Middleman{
		errorHandler: errorHandler,
		httpServer: http.Server{
			Addr:         addr,
			Handler:      http.HandlerFunc(mm.mainHandler),
			TLSNextProto: tlsNextProto,
		},
	}
}

// ListenAndServeTLS starts the https server
func (mm *Middleman) ListenAndServeTLS(certFile, keyFile string) error {
	// Start the listener, and if an error occures, pass it up to the caller
	err :=
		mm.httpServer.ListenAndServeTLS(certFile, keyFile)

	return err
}

// ListenAndServe starts the http server
func (mm *Middleman) ListenAndServe() error {
	err := mm.httpServer.ListenAndServe()

	return err
}

// emitError calls the error handler callback to inform the user of an error
// and returns if execution should continue
func (mm *Middleman) emitError(err error) bool {
	if mm.errorHandler != nil {
		return mm.errorHandler(err)
	}

	return true
}

// mainHandler is the main function that receives all
// requests and calls the correct middlewares
func (mm *Middleman) mainHandler(res http.ResponseWriter, req *http.Request) {
	// Create a store to hold information between middlewares
	store := Store{}

	_, err := mm.runMiddlewares(res, req, store)

	if err != nil {
		log.Println("[mainHandler error]:", err.Error())
	}
}

// addMiddleware adds a middleware to the middleware store
func (mm *Middleman) addMiddleware(path string, method string,
	middleware Middleware) {
	mm.handlers = append(mm.handlers, middlewareHandler{
		middleware,
		path,
		method,
	})
}

// runMiddlewares runs middlewares on a request
// Returns a bool value to indicate if execution stopped
// Returns an error if any occured
func (mm *Middleman) runMiddlewares(res http.ResponseWriter, req *http.Request,
	store Store) (bool, error) {
	// Indication weather execution should be stopped
	cont := true

	// Define the end function
	end := func() {
		cont = false
	}

	// Iterate over all handlers
	for _, handler := range mm.handlers {
		// If the middleware called the end function, middleware execution
		// should be stopped
		if !cont {
			break
		}

		// Match the regex of the handler to the request's uri path
		regexMatch, err := regexp.MatchString("^"+handler.path+"$",
			req.URL.Path)

		if err != nil {
			continueAfterError :=
				mm.emitError(errors.New("[Regex matching error]: " +
					err.Error()))

			return continueAfterError,
				errors.New("[Regex matching error]: " + err.Error())
		}

		if regexMatch && handler.method == req.Method {
			err := handler.middleware(res, req, store, end)

			// If an error occured in the middleware, emit the error
			if err != nil {
				errMsg := "[Method: " + req.Method +
					" Path: " + req.RequestURI + "]: "

				// Raise error emitter and decide to continue or break
				continueAfterError :=
					mm.emitError(errors.New(errMsg + err.Error()))

				// If emitError returns false, break execution
				if !continueAfterError {
					break
				}
			}
		}
	}

	return cont, nil
}
