package caf

import (
	"log"
	"os"

	"github.com/Creespye/caf/internal/pkg/configs"
	"github.com/Creespye/caf/internal/pkg/proxy"
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

	proxy.Start(proxy.Config{
		Addr:   ":" + config.Out.Port,
		Target: config.In.Targets[0].GetURL(),
		Cert:   config.Out.CertificatePath,
		Key:    config.Out.KeyPath,
	})
}
