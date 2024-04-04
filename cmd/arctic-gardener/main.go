package main

import (
  "flag"
  "log"
  "time"

  "github.com/charlesmcchan/arctic-gardener/internal/adc"
  "github.com/charlesmcchan/arctic-gardener/internal/config"
  "github.com/charlesmcchan/arctic-gardener/internal/gpio"
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
    if c.checkLastOn() {
      duration, err := time.ParseDuration(c.Duration)
      if err != nil {
        log.Fatal("Error parsing duration:", err)
      }
      gpio.On(c.Pin.Pump, duration)
      c.UpdateLastOn(*configPath)
    } else {
      log.Println("Last on is too recent. Skipping...")
    }
  } else {
    gpio.Off(c.Pin.Pump)
  }
}
