package configs

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

/*
A struct that represents a JSON object that contains the configuration of our project.
*/
type Configuration struct {
	In                 In  `json:"in"`
	Out                Out `json:"out"`
	LocalRootDirectory string
}

/*
A function that creates a new Configuration struct.
*/
func NewConfiguration() Configuration {
	localRoot, _ := os.Getwd()
	return Configuration{LocalRootDirectory: localRoot}
}

/*
This function get a pointer to a Configuration struct and populates it
with configuration from a JSON file.
*/
func GetConf(config *Configuration) error {
	var err error
	workingDirectory, err := os.Getwd()
	if err != nil {
		return err
	}

	// Open the configuration file.
	file, err := os.Open(workingDirectory + "/settings/config.dev.json")
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

	// Return the error
	return err
}
