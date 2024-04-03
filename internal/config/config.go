package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Threshold struct {
	Wet int `yaml:"wet"`
	Dry int `yaml:"dry"`
}

type Pin struct {
	Sensor int `yaml:"sensor"`
	Pump int `yaml:"pump"`
}

// Config represents the configuration for the application.
type Config struct {
	Pin Pin `yaml:"pin"`
	Threshold Threshold `yaml:"threshold"`
}

// Init reads the configuration from the provided file path
// and returns a Config object.
func NewConfig(configPath string) (Config, error) {
	var config Config

	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(configFile, &config); err != nil {
		return config, err
	}

	if err = checkConfig(config) ; err != nil {
		return config, err
	}

	return config, nil
}

// checkConfig validates the provided configuration.
// If any validation fails, it returns an error.
func checkConfig(c Config) error {
	if c.Threshold.Wet < 0 || c.Threshold.Wet > 100 {
		return fmt.Errorf("WetThreshold must be between 0 and 100")
	}
	if c.Threshold.Dry < 0 || c.Threshold.Dry > 100 {
		return fmt.Errorf("DryThreshold must be between 0 and 100")
	}
	if c.Pin.Sensor == 0 {
		return fmt.Errorf("Pins.Sensor must not be 0")
	}
	if c.Pin.Pump == 0 {
		return fmt.Errorf("Pins.Pump must not be 0")
	}
	return nil
}