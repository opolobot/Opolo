package utils

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var config *Config

const configFileName string = "./config.yml"

// Config represents the bot configuration file.
type Config struct {
	LogChannel string `yaml:"logChannel"`
	Prefix     string `yaml:"prefix"`
	Status     string `yaml:"status"`
	Token      string `yaml:"token"`
}

// GetConfig gets the config from the filesystem followed by a singleton
func GetConfig() *Config {
	if config == nil {
		var err error
		config, err = readConfig()
		if err != nil {
			log.Fatal("Failed to load config: ", err)
		}
	}

	return config
}

// Found here: https://stackoverflow.com/questions/30947534/how-to-read-a-yaml-file
func readConfig() (*Config, error) {
	log.Printf("Loading config file at %v\n", configFileName)
	buf, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("Bad config file: %v", err)
	}

	return c, nil
}
