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
	// Open the configuration file.
	file, err := os.Open(config.SettingsFilePath)
	if err != nil {
		return err
	}

	// When the function is done, close the file.
	defer file.Close()

	// Read the data from the file.
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	// Unmarshal the json bytes into the.
	err = json.Unmarshal(bytes, config)

	// Create settings folder path from setting file path for extracting relative certs path
	SettingsFolderPath := path.Dir(strings.ReplaceAll(config.SettingsFilePath, "\\", "/")) + "/"

	config.Out.CertificatePath =
		SettingsFolderPath + config.Out.CertificatePath

	config.Out.KeyPath = SettingsFolderPath + config.Out.KeyPath

	// Return the error
	return err
}
