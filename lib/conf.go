package lib

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ConfFileName is Location of config file relative to the dir of the executable.
const ConfFileName string = "./config.yml"

// Config represents the bot configuration file.
type Config struct {
	LogChannel string `yaml:"logChannel"`
	Prefix     string `yaml:"prefix"`
	Status     string `yaml:"status"`
	Token      string `yaml:"token"`
}

// FetchConf fetches a config file and populates the structure.
// Found here: https://stackoverflow.com/questions/30947534/how-to-read-a-yaml-file
func FetchConf() (*Config, error) {
	buf, err := ioutil.ReadFile(ConfFileName)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("Bad config file %q: %v", ConfFileName, err)
	}

	return c, nil
}
