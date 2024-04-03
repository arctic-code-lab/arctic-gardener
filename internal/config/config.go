package config

import (
  "fmt"
  "io/ioutil"
  "slices"

  "github.com/charlesmcchan/arctic-gardener/internal/control"
  "gopkg.in/yaml.v2"
)

type Threshold struct {
  Wet int `yaml:"wet"`
  Dry int `yaml:"dry"`
}

type Pin struct {
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
  if c.Threshold.Wet < adc.AdcMinValue || c.Threshold.Wet > adc.AdcMaxValue {
    return fmt.Errorf("WetThreshold must be between 0 and 100")
  }
  if c.Threshold.Dry < adc.AdcMinValue || c.Threshold.Dry > adc.AdcMaxValue {
    return fmt.Errorf("DryThreshold must be between 0 and 100")
  }
  if !slices.Contains(gpioPins, c.Pin.Pump) {
    return fmt.Errorf("%d is not a valid GPIO pin", c.Pin.Pump)
  }
  return nil
}