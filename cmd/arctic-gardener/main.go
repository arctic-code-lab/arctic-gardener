package main

import (
  "flag"
  "log"

  "github.com/charlesmcchan/arctic-gardener/internal/control/adc"
  "github.com/stianeikeland/go-rpio/v4"
)

func main() {
  log.SetFlags(log.LstdFlags | log.Lshortfile)

  configPath := flag.String("c", "configs/configs.yaml", "path to config file")
  flag.Parse()

  c, err := config.NewConfig(*configPath)
  if err != nil {
    log.Fatal("Config error:", err)
  }

  if err := rpio.Open(); err != nil {
    log.Fatal(err)
  }
  defer rpio.Close()

  log.Printf("Config: %+v\n", c)

  pin := rpio.Pin(c.Pin.Sensor)
  pin.Input()
  pin.PullDown()
  sensorResult := pin.Read()
  log.Printf("Pin %d = %d\n", pin, sensorResult)

  pin = rpio.Pin(c.Pin.Pump)
  pin.Output()
  if sensorResult == rpio.Low {
    log.Println("Soil is wet")
    pin.Low()
  } else {
    log.Println("Soil is dry")
    pin.High()
  }

  adc.NewAdc()
}
