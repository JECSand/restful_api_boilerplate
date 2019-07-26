/*
Author: Connor Sanders
MIT License
RESTful API Boilerplate
7/19/2019
*/

package configuration

import (
	"encoding/json"
	"fmt"
	"os"
)

// Configuration is a struct designed to hold the applications variable configuration settings
type Configuration struct {
	MongoURI                   string
	Secret                     string
	Database                   string
	MasterAdminUsername        string
	MasterAdminEmail           string
	MasterAdminInitialPassword string
	DefaultAdminGroup          string
	DefaultDateTimeLayout      string
	HTTPS                      string
	Cert                       string
	Key                        string
}

// ConfigurationSettings is a function that reads a json configuration file and outputs a Configuration struct
func ConfigurationSettings(env string) Configuration {
	confFile := "conf.json"
	if env == "test" {
		confFile = "test_conf.json"
	}
	file, _ := os.Open(confFile)
	decoder := json.NewDecoder(file)
	configurationSettings := Configuration{}
	err := decoder.Decode(&configurationSettings)
	if err != nil {
		fmt.Println("error: ", err)
	}
	return configurationSettings
}
