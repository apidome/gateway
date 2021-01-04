package configs

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"fmt"
)

var config *Configuration

// Configuration is a struct that represents a JSON object
// that contains the configuration of our project.
type Configuration struct {
	General          General `json:"general"`
	In               In      `json:"in"`
	Out              Out     `json:"out"`
	SettingsFilePath string
}

// GetConfiguration function gets a pointer to a Configuration
// struct and populates it with configuration from a JSON file.
func GetConfiguration() (*Configuration, error) {
	// if first time GetConfiguration called
	if config == nil {
		args := os.Args

		if len(args) < 2 {
			log.Panicln("Not enough arguments, probably " +
				"missing configuration file path.")
		}

		config = &Configuration{
			SettingsFilePath: args[1],
		}

		err := readConf(config)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

// readConf reads configurations from a file and stores it in the
// received Configuration pointer.
func readConf(config *Configuration) error {
	bytes, err := ioutil.ReadFile(config.SettingsFilePath)

	// Unmarshal the json bytes into the.
	err = json.Unmarshal(bytes, config)
	if err != nil {
		return err
	}

	// Read the schema of each endpoint from
	// file and set it in the schema field.
	for _, target := range config.In.Targets {
		for _, api := range target.Apis {
			for _, endpoint := range api.Endpoints {
				// Read the data from file.
				schema, err :=
					ioutil.ReadFile(endpoint.Schema)

				if err != nil {
					return err
				}

				// Set the actual schema in the endpoint.
				endpoint.Schema = string(schema)
			}
		}
	}

	config.Out.CertificatePath =filepath.FromSlash(config.Out.CertificatePath)

	config.Out.KeyPath = filepath.FromSlash(config.Out.KeyPath)

	fmt.Println(config.Out.CertificatePath, config.Out.KeyPath)

	// Return the error
	return err
}
