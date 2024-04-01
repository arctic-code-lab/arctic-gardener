package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/charlesmcchan/arctic-gardener/internal/config"
	"github.com/stianeikeland/go-rpio/v4"
)

func main() {
	configPath := flag.String("c", "configs/configs.yaml", "path to config file")
	flag.Parse()

	c, err := config.NewConfig(*configPath)
	if err != nil {
		fmt.Println("Config error:", err)
		os.Exit(1)
	}

	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()

	fmt.Printf("Config: %+v\n", c)

	pin := rpio.Pin(c.Pin.Sensor)
	pin.Input()
	pin.PullDown()
	sensorResult := pin.Read()
	fmt.Printf("Pin %d = %d\n", pin, sensorResult)

	pin = rpio.Pin(c.Pin.Pump)
	pin.Output()
	if sensorResult == rpio.Low {
		fmt.Println("Soil is wet")
		pin.Low()
	} else {
		fmt.Println("Soil is dry")
		pin.High()
	}
}