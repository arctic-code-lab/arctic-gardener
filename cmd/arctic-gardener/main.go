package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/charlesmcchan/arctic-gardener/internal/adc"
	"github.com/charlesmcchan/arctic-gardener/internal/config"
	"github.com/charlesmcchan/arctic-gardener/internal/gpio"
)

const (
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	reset  = "\033[m"
)

func main() {
	flagConfig := flag.String("c", "configs.yaml", "path to config file")
	flagDryRun := flag.Bool("d", false, "dry run")
	flagSimple := flag.Bool("s", false, "print humidity and exit")
	flag.Parse()

	c, err := config.NewConfig(*flagConfig)
	if err != nil {
		log.Fatal("Error reading config:", err)
	}

	duration, err := time.ParseDuration(c.Duration)
	if err != nil {
		log.Fatal("Error parsing duration:", err)
	}

	// Reset the state of the pump
	gpio.Off(c.Pin.Pump)

	adcReading := adc.Read(c.Pin.Sensor)
	humidity := 100 * (c.Threshold.Dry - adcReading) / (c.Threshold.Dry - c.Threshold.Wet)
	threshold := c.Threshold.Dry - c.Threshold.Percent*(c.Threshold.Dry-c.Threshold.Wet)/100

	if *flagSimple {
		fmt.Println(humidity)
		return
	}
	log.Printf("%d%% %d | %d%% %d\n", humidity, adcReading, c.Threshold.Percent, threshold)

	if adcReading < threshold || *flagDryRun {
		return
	}
	if !c.ShoudRun() {
		log.Printf("%sLast run is too recent. Skipping...%s", yellow, reset)
		return
	}

	log.Printf("%sTurning on for %s...%s", green, duration, reset)
	gpio.On(c.Pin.Pump, duration)

	log.Printf("Updating lastRun...\n")
	c.UpdateLastRun(*flagConfig)
}
