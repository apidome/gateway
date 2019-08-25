package caf

import (
	"log"
	"os"

	"github.com/Creespye/caf/internal/pkg/middleman"
	"github.com/Creespye/caf/internal/pkg/proxy"
	"github.com/Creespye/caf/internal/pkg/proxymiddlewares"

	"github.com/Creespye/caf/internal/pkg/configs"
)

// Start starts CAF
func Start() {
	args := os.Args

	if len(args) < 2 {
		log.Panicln("Not enough arguments, probably missing configuration file path.")
	}

	settingsFolder := args[1]

	// Create a new configuration struct.
	config := configs.NewConfiguration(settingsFolder)

	// Populate the configuration struct.
	err := configs.GetConf(&config)
	if err != nil {
		log.Panicln("Could not load configuration correctly:", err)
		os.Exit(2)
	}

	// Allocate a new middleman struct for the proxy
	var mmProxy middleman.Middleman

	// Create a new middleman webserver for the proxy
	middleman.NewMiddleman(&mmProxy,
		":"+config.Out.Port,
		func(err error) bool {
			log.Println("[Middleman Proxy Error]: " + err.Error())

			return false
		})

	// Create a new proxy server
	pr := proxy.NewProxy()

	urlWithNoProto :=
		config.In.Targets[0].Host + ":" + config.In.Targets[0].Port

	urlWithNoProto = "localhost:30000"

	mmProxy.Connect("", proxymiddlewares.PrintConnections())

	mmProxy.Connect("",
		proxymiddlewares.TunnelConnection(&pr, urlWithNoProto))

	// Allocate a new midleman struct for the intercepter
	var mmIntercepter middleman.Middleman

	// Create a middleman webserver
	middleman.NewMiddleman(&mmIntercepter,
		":"+config.Out.Port+"0",
		func(err error) bool {
			log.Println("[Middleman Error]: " + err.Error())

			return false
		})

	// ================ Web server request handling begins here ===============

	mmIntercepter.All("*", middleman.RouteLogger())

	mmIntercepter.All("/.*", middleman.BodyReader())

	// ==================== Request proxy code begins here ====================

	// ===================== Request proxy code ends here =====================

	mmIntercepter.All("/.*",
		proxymiddlewares.CreateTargetRequest(&pr,
			config.In.Targets[0].GetURL()))

	mmIntercepter.All("/.*", proxymiddlewares.SendTargetRequest(&pr))

	mmIntercepter.All("/.*", proxymiddlewares.ReadTargetResponseBody(&pr))

	// =================== Response proxy code begins here ====================

	// =================== Response proxy code ends here ======================

	mmIntercepter.All("/.*", proxymiddlewares.SendTargetResponse(&pr))

	// ================ Web server request handling ends here =================

	go func() {
		log.Println("[Proxy is listening on]:", config.Out.Port)
		err := mmProxy.ListenAndServe()

		if err != nil {
			log.Fatal("[Failed creating the proxy server]:", err)
		}
	}()

	log.Println("[Intercepter is listening on]:", config.Out.Port+"0")
	err = mmIntercepter.ListenAndServeTLS(config.Out.CertificatePath,
		config.Out.KeyPath)

	// If an error occured, print a message
	if err != nil {
		log.Fatalln("[Failed creating the intercepter server]:", err)
	}
}
