package caf

import (
	"log"
	"net/http"

	"github.com/Creespye/caf/internal/pkg/configs"
	"github.com/Creespye/caf/internal/pkg/middleman"
	"github.com/Creespye/caf/internal/pkg/proxy"
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
