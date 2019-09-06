package caf

import (
	"log"
	"net/http"

	"github.com/Creespye/caf/internal/pkg/configs"
	"github.com/Creespye/caf/internal/pkg/middleman"
	"github.com/Creespye/caf/internal/pkg/proxy"
	"github.com/Creespye/caf/internal/pkg/proxymiddlewares"
	"github.com/Creespye/caf/internal/pkg/validators"
	"github.com/Creespye/caf/internal/pkg/validators/jsonvalidator"
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

	var reverseProxy middleman.Middleman

	initReverseProxy(&reverseProxy,
		config.Out.Port,
		config.In.Targets[0].GetURL())

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
	AddValidationMiddlewares(reverseProxy, config.In.Targets)
	reverseProxy.All("/.*", proxymiddlewares.CreateRequest(pr))
	reverseProxy.All("/.*", proxymiddlewares.SendRequest(pr))
}

func responseProxying(reverseProxy *middleman.Middleman, pr *proxy.Proxy) {
	reverseProxy.All("/.*", proxymiddlewares.ReadResponseBody())
	reverseProxy.All("/.*", proxymiddlewares.SendResponse())
}

func initReverseProxy(reverseProxy *middleman.Middleman,
	listeningPort,
	target string) {
	// Middleman is the underlying webserver/middleware manager for our reverse proxy
	middleman.NewMiddleman(reverseProxy,
		":"+listeningPort,
		middlewareErrorHandler)

	reverseProxy.All("*", middleman.RouteLogger())

	// Read the request body and store it in store["reqeustBody"]
	// for all middlewares to use
	reverseProxy.All("/.*", middleman.BodyReader())

	pr := proxy.NewProxy(target)

	requestProxying(reverseProxy, &pr)
	responseProxying(reverseProxy, &pr)
}

func middlewareErrorHandler(path, method string, err error) bool {
	log.Println("[Middleman Error]:", "\n",
		"[Path]:", path, "\n",
		"[Method]:", method)

	return false
}

// AddValidationMiddlewares gets a reference to a Middleman and a slice of targets
// and creates a new middleware for each endpoint in the targets' apis.
func AddValidationMiddlewares(mm *middleman.Middleman, targets []configs.Target) error {
	// Loop over the targets slice
	for _, target := range targets {
		// For each target loop over its apis
		for _, api := range target.Apis {
			var validator validators.Validator

			// Each api has a validator that filter the api's traffic.
			// Here we decide which validator to create according to the api's type.
			switch api.Type {
			case configs.TypeRest:
				validator = jsonvalidator.NewJsonValidator()
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
