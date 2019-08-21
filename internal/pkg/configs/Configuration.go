package configs

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

const (
	configurationFileName = "config.dev.json"
)

/*
Configuration is a struct that represents a JSON object
that contains the configuration of our project.
*/
type Configuration struct {
	In                 In  `json:"in"`
	Out                Out `json:"out"`
	SettingsFolderPath string
}

/*
NewConfiguration function creates a new Configuration struct.
*/
func NewConfiguration(settingsPath string) Configuration {
	return Configuration{SettingsFolderPath: settingsPath}
}

/*
GetConf function gets a pointer to a Configuration struct and populates it
with configuration from a JSON file.
*/
func GetConf(config *Configuration) error {
	// Open the configuration file.
	file, err := os.Open(config.SettingsFolderPath + configurationFileName)
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

	if !strings.HasSuffix(config.SettingsFolderPath, "/") {
		config.SettingsFolderPath += "/"
	}

	config.Out.CertificatePath =
		config.SettingsFolderPath + config.Out.CertificatePath

	config.Out.KeyPath = config.SettingsFolderPath + config.Out.KeyPath

	// Return the error
	return err
}
