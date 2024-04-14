package config

import (
	"fmt"
	"log"
	"os"
	"slices"
	"time"

	"github.com/charlesmcchan/arctic-gardener/internal/adc"
	"github.com/charlesmcchan/arctic-gardener/internal/gpio"
	"gopkg.in/yaml.v2"
)

type Threshold struct {
	Wet     int `yaml:"wet"`
	Dry     int `yaml:"dry"`
	Percent int `yaml:"percent"`
}

type Pin struct {
	Sensor int `yaml:"sensor"`
	Pump   int `yaml:"pump"`
}

type Config struct {
	Pin       Pin       `yaml:"pin"`
	Threshold Threshold `yaml:"threshold"`
	Duration  string    `yaml:"duration"`
	Interval  string    `yaml:"interval"`
	LastRun   string    `yaml:"lastRun"`
}

// Init reads the configuration from the provided file path
// and returns a Config object.
func NewConfig(configPath string) (Config, error) {
	var c Config
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return c, err
	}
	if err := yaml.Unmarshal(configFile, &c); err != nil {
		return c, err
	}
	if err = c.checkConfig(); err != nil {
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
	if err := os.WriteFile(configPath, configYaml, 0644); err != nil {
		return err
	}
	return nil
}

// CheckLastRun checks if the last on time is older than the interval duration.
func (config Config) CheckLastRun() bool {
	if config.LastRun == "" {
		log.Println("Last run undefined. Proceeding...")
		return true
	}

	now := time.Now()
	lastOnTime, err := time.Parse(time.RFC3339, config.LastRun)
	if err != nil {
		log.Fatal("Error parsing lastOn time:", err)
	}

	interval, err := time.ParseDuration(config.Interval)
	if err != nil {
		log.Fatal("Error parsing interval:", err)
	}
	timeDifference := now.Sub(lastOnTime)

	return timeDifference >= interval
}

// UpdateLastRun updates the LastRun field in the provided configuration
func (config Config) UpdateLastRun(configPath string) {
	now := time.Now()

	config.LastRun = now.Format(time.RFC3339)
	err := UpdateConfig(configPath, config)
	if err != nil {
		log.Fatal("Error writing config:", err)
	}
}

// checkConfig validates the provided configuration.
// If any validation fails, it returns an error.
func (config Config) checkConfig() error {
	if config.Threshold.Wet < adc.AdcMinValue || config.Threshold.Wet > adc.AdcMaxValue {
		return fmt.Errorf("WetThreshold must be between %d and %d", adc.AdcMinValue, adc.AdcMaxValue)
	}
	if config.Threshold.Dry < adc.AdcMinValue || config.Threshold.Dry > adc.AdcMaxValue {
		return fmt.Errorf("DryThreshold must be between %d and %d", adc.AdcMinValue, adc.AdcMaxValue)
	}
	if !slices.Contains(gpio.GpioPins, config.Pin.Pump) {
		return fmt.Errorf("%d is not a valid GPIO pin", config.Pin.Pump)
	}
	return nil
}
