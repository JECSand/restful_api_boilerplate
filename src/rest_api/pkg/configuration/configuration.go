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

type Configuration struct {
	MongoURI                       string
	Secret                         string
	Database                       string
	MasterAdminUsername            string
	MasterAdminEmail               string
	MasterAdminInitialPassword     string
	DefaultAdminGroup              string
	DefaultDateTimeLayout          string
	HTTPS                          string
	Cert                           string
	Key                            string
}

func ConfigurationSettings() Configuration {
	file, _ := os.Open("conf.json")
	//defer file.Close()
	decoder := json.NewDecoder(file)
	configurationSettings := Configuration{}
	err := decoder.Decode(&configurationSettings)
	if err != nil {
		fmt.Println("error: ", err)
	}
	return configurationSettings
}
