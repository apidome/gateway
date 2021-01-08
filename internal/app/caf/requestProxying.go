package caf

import (
	"github.com/apidome/gateway/internal/pkg/middleman"
	"github.com/apidome/gateway/internal/pkg/proxy"
)

// requestProxying assembles all client request middlewares
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

	addValidationMiddlewares(reverseProxy, config.In.Targets)
}
