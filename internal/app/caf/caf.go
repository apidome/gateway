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

	// Create a middleman webserver
	mm := middleman.NewMiddleman(middleman.Config{
		Addr:     ":" + config.Out.Port,
		CertFile: config.Out.CertificatePath,
		KeyFile:  config.Out.KeyPath,
	}, func(err error) bool {
		log.Println("[Middleman Error]: " + err.Error())

		return false
	})

	// Create a new proxy server
	proxy := proxy.NewProxy()

	// ================ Web server request handling begins here ===============

	mm.All("/.*", middleman.RouteLogger())

	mm.All("/.*", middleman.BodyReader())

	// ==================== Request proxy code begins here ====================

	// ===================== Request proxy code ends here =====================

	mm.All("/.*",
		proxymiddlewares.CreateTargetRequest(&proxy,
			config.In.Targets[0].GetURL()))

	mm.All("/.*", proxymiddlewares.SendTargetRequest(&proxy))

	mm.All("/.*", proxymiddlewares.ReadTargetResponseBody(&proxy))

	// =================== Response proxy code begins here ====================

	// =================== Response proxy code ends here ======================

	mm.All("/.*", proxymiddlewares.SendTargetResponse(&proxy))

	// ================ Web server request handling ends here =================

	log.Println("[Middleman is listening on]:", config.Out.Port)
	log.Println("[Proxy is fowrarding to]:", config.In.Targets[0].GetURL())

	// Begin listening
	err = mm.ListenAndServeTLS()

	// If an error occured, print a message
	if err != nil {
		log.Fatalln("[Failed creating a server]:", err)
	}
}
