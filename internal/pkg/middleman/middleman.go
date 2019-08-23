package middleman

import (
	"errors"
	"log"
	"net/http"
	"regexp"
)

// Config holds the configurations for the underlying web server
type Config struct {
	Addr     string
	CertFile string
	KeyFile  string
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
	store *Store, end End) error

// handler is a struct that hold middleware information
type handler struct {
	middleware Middleware
	path       string
	method     string
}

// Middleman is a struct that holds all middlewares
type Middleman struct {
	config       Config
	handlers     []handler
	errorHandler func(error)
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
func NewMiddleman(config Config, errorHandler func(error)) Middleman {
	return Middleman{
		config:       config,
		errorHandler: errorHandler,
	}
}

// ListenAndServeTLS starts the https server
func (mm *Middleman) ListenAndServeTLS() error {
	http.HandleFunc("/", mm.mainHandler)

	// Start the listener, and if an error occures, pass it up to the caller
	err := http.ListenAndServeTLS(mm.config.Addr, mm.config.CertFile, mm.config.KeyFile, nil)

	return err
}

// ListenAndServe starts the http server
func (mm *Middleman) ListenAndServe() error {
	http.HandleFunc("/", mm.mainHandler)

	err := http.ListenAndServe(mm.config.Addr, nil)

	return err
}

// emitError calls the error handler callback to inform the user of an error
func (mm *Middleman) emitError(err error) {
	if mm.errorHandler != nil {
		mm.errorHandler(err)
	}
}

// mainHandler is the main function that receives all
// requests and calls the correct middlewares
func (mm *Middleman) mainHandler(res http.ResponseWriter, req *http.Request) {
	// Create a store to hold information between middlewares
	store := Store{}

	_, err := mm.runMiddlewares(res, req, &store)

	if err != nil {
		log.Println("[mainHandler error]:", err.Error())
	}

	res.Write([]byte{0})
}

// addMiddleware adds a middleware to the middleware store
func (mm *Middleman) addMiddleware(path string, method string,
	middleware Middleware) {
	mm.handlers = append(mm.handlers, handler{
		middleware,
		path,
		method,
	})
}

// runMiddlewares runs middlewares on a request
func (mm *Middleman) runMiddlewares(res http.ResponseWriter, req *http.Request,
	store *Store) (bool, error) {
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
			req.RequestURI)

		if err != nil {
			mm.emitError(errors.New("[Regex matching error]: " + err.Error()))

			return false,
				errors.New("[Regex matching error]: " + err.Error())
		}

		if regexMatch && handler.method == req.Method {
			err := handler.middleware(res, req, store, end)

			// If an error occured in the middleware, emit the error
			if err != nil {
				errMsg := "[Method: " + req.Method +
					" Path: " + req.RequestURI + "]: "

				mm.emitError(errors.New(errMsg + err.Error()))

				// Break middleware execution when an error occured
				break
			}
		}
	}

	return cont, nil
}
