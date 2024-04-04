package config

import (
  "fmt"
  "io/ioutil"
  "log"
  "slices"
  "time"

  "github.com/charlesmcchan/arctic-gardener/internal/adc"
  "github.com/charlesmcchan/arctic-gardener/internal/gpio"
  "gopkg.in/yaml.v2"
)

type Threshold struct {
  Wet int `yaml:"wet"`
  Dry int `yaml:"dry"`
  Percent int `yaml:"percent"`
}

type Pin struct {
  Sensor int `yaml:"sensor"`
  Pump int `yaml:"pump"`
}

// Config represents the configuration for the application.
type Config struct {
  Pin Pin `yaml:"pin"`
  Threshold Threshold `yaml:"threshold"`
  Duration string `yaml:"duration"`
  Interval string `yaml:"interval"`
  LastOn string `yaml:"lastOn"`
}

// Init reads the configuration from the provided file path
// and returns a Config object.
func NewConfig(configPath string) (Config, error) {
  var c Config
  configFile, err := ioutil.ReadFile(configPath)
  if err != nil {
    return c, err
  }
  if err := yaml.Unmarshal(configFile, &c); err != nil {
    return c, err
  }
  if err = c.checkConfig() ; err != nil {
    return c, err
  }
  return c, nil
}

// UpdateConfig writes the provided configuration to the provided file path.
func UpdateConfig(configPath string, config Config) error {
  configYaml, err := yaml.Marshal(config)
  if err != nil {
    return err
  }
  if err := ioutil.WriteFile(configPath, configYaml, 0644); err != nil {
    return err
  }
  return nil
}

func (config Config) CheckLastOn() bool {
  if config.LastOn == "" {
    log.Println("Last on undefined. Proceeding...")
    return true
  }

  now := time.Now()
  lastOnTime, err := time.Parse(time.RFC3339, config.LastOn)
  if err != nil {
    log.Fatal("Error parsing lastOn time:", err)
  }
  log.Printf("Now: %s\n", now.Format(time.RFC3339))
  log.Printf("Last on: %s\n", lastOnTime.Format(time.RFC3339))

  // Compare the difference between current time and last on time
  // with a threshold duration, e.g. 24 hours
  interval, err := time.ParseDuration(config.Interval)
  if err != nil {
    log.Fatal("Error parsing interval:", err)
  }
  timeDifference := now.Sub(lastOnTime)

  return timeDifference >= interval
}

func (config Config) UpdateLastOn(configPath string) {
  now := time.Now()

  config.LastOn = now.Format(time.RFC3339)
  err := UpdateConfig(configPath, config)
  if err != nil {
    log.Fatal("Error writing config:", err)
  }

  log.Printf("Updating Last On: %s\n", config.LastOn)
}

// checkConfig validates the provided configuration.
// If any validation fails, it returns an error.
func (c Config) checkConfig() error {
  if c.Threshold.Wet < adc.AdcMinValue || c.Threshold.Wet > adc.AdcMaxValue {
    return fmt.Errorf("WetThreshold must be between %d and %d", adc.AdcMinValue, adc.AdcMaxValue)
  }
  if c.Threshold.Dry < adc.AdcMinValue || c.Threshold.Dry > adc.AdcMaxValue {
    return fmt.Errorf("DryThreshold must be between %d and %d", adc.AdcMinValue, adc.AdcMaxValue)
  }
  if !slices.Contains(gpio.GpioPins, c.Pin.Pump) {
    return fmt.Errorf("%d is not a valid GPIO pin", c.Pin.Pump)
  }
  return nil
}