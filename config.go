package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ArubaEndpoint string `yaml:"arubaEndpoint"`
	ArubaTokens   []struct {
		ArubaAccessToken  string `yaml:"arubaAccessToken"`
		ArubaRefreshToken string `yaml:"arubaRefreshToken"`
	} `yaml:"arubaTokens"`
	ArubaApplicationCredentials []struct {
		ClientID     string `yaml:"clientId"`
		ClientSecret string `yaml:"clientSecret"`
	} `yaml:"arubaApplicationCredentials"`
	ExporterConfig []struct {
		ExporterEndpoint string `yaml:"exporterEndpoint"`
		ExporterPort     string `yaml:"exporterPort"`
	} `yaml:"exporterConfig"`
}

func readConfig(c *Config, configPath string) {
	// Read the YAML file
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the YAML data into a Config struct
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Reading confugration from file: ", configPath)

}
