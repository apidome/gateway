package caf

import (
	"github.com/Creespye/caf/internal/pkg/configs"
	"github.com/Creespye/caf/internal/pkg/middleman"
	"github.com/Creespye/caf/internal/pkg/proxy"
	"github.com/Creespye/caf/internal/pkg/proxymiddlewares"
	"github.com/Creespye/caf/internal/pkg/validators"
	"github.com/Creespye/caf/internal/pkg/validators/jsonvalidator"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
)

var config *configs.Configuration

// Start starts CAF
func Start() {
	var err error

	// Initialize and Populate the configuration struct.
	config, err = configs.GetConfiguration()
	if err != nil {
		log.Panicln("Could not load configuration correctly:", err)
	}

	var prx proxy.Proxy

	proxy.InitProxy(&prx, config.In.Targets[0].GetURL())

	var reverseProxy middleman.Middleman

	middleman.InitMiddleman(&reverseProxy,
		":"+config.Out.Port,
		middlewareErrorHandler)

	requestProxying(&reverseProxy, &prx)

	responseProxying(&reverseProxy, &prx)

	reverseProxy.All("/.*", defaultMiddleware())

	log.Println("[Reverse proxy is listening on]:", config.Out.Port)

	if config.Out.SSL {
		err = reverseProxy.ListenAndServeTLS(config.Out.CertificatePath,
			config.Out.KeyPath)
	} else {
		err = reverseProxy.ListenAndServe()
	}

	// If an error occured, print a message
	if err != nil {
		log.Fatalln("[Reverse proxy set up failed]:", err)
	}
}

func requestProxying(reverseProxy *middleman.Middleman, pr *proxy.Proxy) {
	// Log all incoming requests' routes
	reverseProxy.All("/.*", middleman.RouteLogger())

	// Read variables from the request path and parameters from the
	// request query
	reverseProxy.All("/.*", middleman.VariablesReader())
	reverseProxy.All("/.*", middleman.ParametersReader())

	// Read the request body and store it in store["reqeustBody"]
	// for all middlewares to use
	reverseProxy.All("/.*", middleman.BodyReader())

	AddValidationMiddlewares(reverseProxy, config.In.Targets)

	reverseProxy.All("/.*", proxymiddlewares.CreateRequest(pr))
}

func responseProxying(reverseProxy *middleman.Middleman, pr *proxy.Proxy) {
	reverseProxy.All("/.*", proxymiddlewares.SendRequest(pr))
	reverseProxy.All("/.*", proxymiddlewares.ReadResponseBody())
	reverseProxy.All("/.*", proxymiddlewares.SendResponse())
}

func middlewareErrorHandler(path, method string, err error) bool {
	log.Println("[Middleman Error]:", err.Error(), "\n",
		"[Path]:", path, "\n",
		"[Method]:", method)

	return false
}

func defaultMiddleware() middleman.Middleware {
	return func(res http.ResponseWriter, req *http.Request,
		store middleman.Store, end middleman.End) error {
		res.WriteHeader(404)
		return nil
	}
}

// AddValidationMiddlewares gets a reference to a Middleman and a slice of targets
// and creates a new middleware for each endpoint in the targets' apis.
func AddValidationMiddlewares(mm *middleman.Middleman, targets []configs.Target) error {
	// Loop over the targets slice
	for _, target := range targets {
		// For each target loop over its apis
		for index, api := range target.Apis {
			var err error
			var validator validators.Validator

			// Each api has a validator that filter the api's traffic.
			// Here we decide which validator to create according to the api's type.
			switch api.Type {
			case configs.TypeRest:
				validator, err = jsonvalidator.NewJsonValidator(api.Version)
				if err != nil {
					return errors.Wrap(err, "failed to created validator for number - "+strconv.Itoa(index))
				}
			default:
				log.Print("[Proxy WARNING]: Invalid API Type - " + api.Type)
			}

			// For each api loop over its endpoints
			for _, endpoint := range api.Endpoints {
				//Add the endpoint's schema to the api's validator.
				err := validator.LoadSchema(endpoint.Path, endpoint.Method, []byte(endpoint.Schema))
				if err != nil {
					log.Print("[Proxy ERROR]: Failed to load schema for endpoint - " + endpoint.Path + ", Error: " + err.Error())
					return err
				}

				// Creating a new ValidateRequest middleware with the appropriate HTTP method.
				switch endpoint.Method {
				case http.MethodGet:
					mm.Get(endpoint.Path, proxymiddlewares.ValidateRequest(endpoint.Path,
						endpoint.Method,
						validator))
				case http.MethodPost:
					mm.Post(endpoint.Path, proxymiddlewares.ValidateRequest(endpoint.Path,
						endpoint.Method,
						validator))
				case http.MethodPut:
					mm.Put(endpoint.Path, proxymiddlewares.ValidateRequest(endpoint.Path,
						endpoint.Method,
						validator))
				case http.MethodDelete:
					mm.Delete(endpoint.Path, proxymiddlewares.ValidateRequest(endpoint.Path,
						endpoint.Method,
						validator))
				case "ALL":
					mm.All(endpoint.Path, proxymiddlewares.ValidateRequest(endpoint.Path,
						endpoint.Method,
						validator))
				default:
					log.Print("[Proxy WARNING]: Invalid method - " + endpoint.Method + " for endpoint - " + endpoint.Path)
				}

				log.Print("[Proxy DEBUG]: Added middleware for - " + endpoint.Method + " " + endpoint.Path)
			}
		}
	}

	return nil
}
