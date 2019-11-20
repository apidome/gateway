package caf

import (
	"github.com/Creespye/caf/internal/pkg/middleman"
	"github.com/Creespye/caf/internal/pkg/proxy"
	"github.com/Creespye/caf/internal/pkg/proxymiddlewares"
)

// responseProxying assembles all target response middlewares
func responseProxying(reverseProxy *middleman.Middleman, pr *proxy.Proxy) {
	reverseProxy.All("/.*", proxymiddlewares.CreateRequest(pr))
	reverseProxy.All("/.*", proxymiddlewares.SendRequest(pr))
	reverseProxy.All("/.*", proxymiddlewares.ReadResponseBody())
	reverseProxy.All("/.*", proxymiddlewares.SendResponse())
}
