package main

import (
  "flag"
  "log"
  "time"

  "github.com/charlesmcchan/arctic-gardener/internal/adc"
  "github.com/charlesmcchan/arctic-gardener/internal/config"
  "github.com/charlesmcchan/arctic-gardener/internal/gpio"
)

const (
	red = "\033[31m"
	green = "\033[32m"
  yellow = "\033[33m"
	reset = "\033[m"
)

func main() {
  log.SetFlags(log.LstdFlags | log.Lshortfile)

  configPath := flag.String("c", "configs/configs.yaml", "path to config file")
  flag.Parse()

  c, err := config.NewConfig(*configPath)
  if err != nil {
    log.Fatal("Error reading config:", err)
  }
  log.Printf("Config: %+v\n", c)

  adcReading := adc.Read(c.Pin.Sensor)
  threshold := c.Threshold.Wet + (c.Threshold.Dry - c.Threshold.Wet) * c.Threshold.Percent / 100
  log.Printf("Reading: %d\n", adcReading)
  log.Printf("Threshold: %d\n", threshold)

  if adcReading >= threshold {
    if c.CheckLastOn() {
      duration, err := time.ParseDuration(c.Duration)
      if err != nil {
        log.Fatal("Error parsing duration:", err)
      }
      log.Printf("%sTurning on pin %d for %s...%s", green, c.Pin.Pump, duration, reset)
      gpio.On(c.Pin.Pump, duration)

      log.Printf("Updating Last On: %s\n", c.LastOn)
      c.UpdateLastOn(*configPath)
    } else {
      log.Printf("%sLast run is too recent. Skipping...%s", yellow, reset)
    }
  } else {
    log.Printf("%sTurning off pin %d%s", red, c.Pin.Pump, reset)
    gpio.Off(c.Pin.Pump)
  }
}
