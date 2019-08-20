package caf

import (
	"fmt"
	"github.com/Creespye/caf/configs/json_configs"
	"github.com/Creespye/caf/internal/pkg/proxy"
	"os"
)

// Start starts CAF
func Start() {
	var err error

	// Create a new configuration struct.
	config := json_configs.NewConfiguration()

	// Populate the configuration struct.
	err = json_configs.GetConf(&config)
	if err != nil {
		fmt.Println("Could not load configuration correctly:", err)
		os.Exit(2)
	}

	proxy.Start(proxy.Config{
		Addr:   ":" + config.Out.Port,
		Target: config.In.Targets[0].GetURL(),
		Cert:   config.LocalRootDirectory + "/settings/certs/localhost/localhost.cert",
		Key:    config.LocalRootDirectory + "/settings/certs/localhost/localhost.key",
	})
}
