package configs

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

var config *Configuration

/*
Configuration is a struct that represents a JSON object
that contains the configuration of our project.
*/
type Configuration struct {
	In               In  `json:"in"`
	Out              Out `json:"out"`
	SettingsFilePath string
}

/*
GetConf function gets a pointer to a Configuration struct and populates it
with configuration from a JSON file.
*/
func GetConfiguration() (*Configuration, error) {
	// if first time GetConfiguration called
	if config == nil {
		args := os.Args

		if len(args) < 2 {
			log.Panicln("Not enough arguments, probably missing configuration file path.")
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

func readConf(config *Configuration) error {
	bytes, err := readDataFromFile(config.SettingsFilePath)

	// Unmarshal the json bytes into the.
	err = json.Unmarshal(bytes, config)
	if err != nil {
		return err
	}

	// Create settings folder path from setting file path for extracting relative certs path
	SettingsFolderPath := path.Dir(strings.ReplaceAll(config.SettingsFilePath, "\\", "/")) + "/"

	// Read the schema of each endpoint from file and set it in the schema field.
	for _, target := range config.In.Targets {
		for _, api := range target.Apis {
			for _, endpoint := range api.Endpoints {
				// Read the data from file.
				schema, err := readDataFromFile(SettingsFolderPath + endpoint.Schema)
				if err != nil {
					return err
				}

				// Set the actual schema in the endpoint.
				endpoint.Schema = string(schema)
			}
		}
	}

	config.Out.CertificatePath =
		SettingsFolderPath + config.Out.CertificatePath

	config.Out.KeyPath = SettingsFolderPath + config.Out.KeyPath

	// Return the error
	return err
}

func readDataFromFile(filePath string) ([]byte, error) {
	// Open the configuration file.
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	// Read the data from the file.
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
