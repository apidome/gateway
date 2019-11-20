package caf

import (
	"github.com/omeryahud/caf/internal/pkg/middleman"
	"github.com/omeryahud/caf/internal/pkg/proxy"
	"github.com/omeryahud/caf/internal/pkg/proxymiddlewares"
)

// responseProxying assembles all target response middlewares
func responseProxying(reverseProxy *middleman.Middleman, pr *proxy.Proxy) {
	reverseProxy.All("/.*", proxymiddlewares.CreateRequest(pr))
	reverseProxy.All("/.*", proxymiddlewares.SendRequest(pr))
	reverseProxy.All("/.*", proxymiddlewares.ReadResponseBody())
	reverseProxy.All("/.*", proxymiddlewares.SendResponse())
}
