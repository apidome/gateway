package caf

import (
	"log"
	"os"

	"github.com/Creespye/caf/internal/pkg/configs"
	"github.com/Creespye/caf/internal/pkg/middleman"
	"github.com/Creespye/caf/internal/pkg/proxy"
	"github.com/Creespye/caf/internal/pkg/proxymiddlewares"
)

// Start starts CAF
func Start() {
	args := os.Args

	if len(args) < 2 {
		log.Panicln("Not enough arguments, probably missing configuration file path.")
	}

	settingsFolder := args[1]

	config := configs.NewConfiguration(settingsFolder)

	err := configs.GetConf(&config)
	if err != nil {
		log.Panicln("Could not load configuration correctly:", err)
		os.Exit(2)
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
		func(err error) bool {
			log.Println("[Middleman Error]: " + err.Error())

			return false
		})

	reverseProxy.All("*", middleman.RouteLogger())

	// Read the request body and store it in store["reqeustBody"]
	// for all middlewares to use
	reverseProxy.All("/.*", middleman.BodyReader())

	pr := proxy.NewProxy(target)

	requestProxying(reverseProxy, &pr)
	responseProxying(reverseProxy, &pr)
}
