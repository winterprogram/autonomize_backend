package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"test/test_app/app/constants"

	"github.com/spf13/viper"
)

var config *viper.Viper

//Init :
func Init(service, env string) {
	body, err := fetchConfiguration(service, env)
	if err != nil {
		fmt.Println("Couldn't load configuration, cannot start. Terminating. Error: " + err.Error())
	}
	parseConfiguration(body)
}

// Make HTTP request to fetch configuration from config server
func fetchConfiguration(service, env string) ([]byte, error) {
	var bodyBytes []byte
	var err error
	if env == constants.DevEnvironment {
		//panic("Couldn't load configuration, cannot start. Terminating. Error: " + err.Error())
		bodyBytes, err = ioutil.ReadFile("test_app/config/config.json")
		if err != nil {
			fmt.Println("Couldn't read local configuration file.", err)
		} else {
			log.Print("using local config.")
		}
	}
	return bodyBytes, err
}

func getEnvOrDefault(envKey, defaultValue string) string {
	var envValue string
	var ok bool
	if envValue, ok = os.LookupEnv(envKey); !ok {
		envValue = defaultValue
	}
	return envValue
}

// Pass JSON bytes into struct and then into Viper
func parseConfiguration(body []byte) {
	var cloudConfig springCloudConfig
	err := json.Unmarshal(body, &cloudConfig)
	if err != nil {
		fmt.Println("Cannot parse configuration, message: " + err.Error())
	}
	for key, value := range cloudConfig.PropertySources[0].Source {
		viper.Set(key, value)
		fmt.Printf("Loading config property > %s - %s \n", key, value)
	}
	if viper.IsSet("server_name") {
		fmt.Println("Successfully loaded configuration for service\n", viper.GetString("server_name"))
	}
}

// Structs having same structure as response from Spring Cloud Config
type springCloudConfig struct {
	Name            string           `json:"name"`
	Profiles        []string         `json:"profiles"`
	Label           string           `json:"label"`
	Version         string           `json:"version"`
	PropertySources []propertySource `json:"propertySources"`
}
type propertySource struct {
	Name   string                 `json:"name"`
	Source map[string]interface{} `json:"source"`
}
