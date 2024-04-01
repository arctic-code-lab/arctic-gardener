package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Pins []int `yaml:"pins"`
	PollIntervalMs int `yaml:"pollIntervalMs"`
	WetThreshold int `yaml:"wetThreshold"`
	DryThreshold int `yaml:"dryThreshold"`
}

func main() {
	// Parse command-line arguments
	configPath := flag.String("c", "configs/configs.yaml", "path to config file")
	flag.Parse()

	// Read config from file
	configFile, err := ioutil.ReadFile(*configPath)
	if err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)
		os.Exit(1)
	}

	// Parse config
	var configs Config
	if err := yaml.Unmarshal(configFile, &configs); err != nil {
		fmt.Println("Failed to parse config:", err)
		os.Exit(1)
	}
	if len(configs.Pins) == 0 {
		fmt.Println("No pins specified in config")
		os.Exit(1)
	}
	fmt.Printf("Config: %+v\n", configs)

	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Signal handling
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <- ch
		fmt.Println("Received signal:", sig)
		rpio.Close()
		os.Exit(0)
	}()

	for {
		fmt.Println("---")
		for _, pin := range configs.Pins {
			pin := rpio.Pin(pin)
			pin.Input()
			pin.PullDown()
			fmt.Printf("Pin %d = %d\n", pin, pin.Read())
		}
		time.Sleep(time.Duration(configs.PollIntervalMs) * time.Millisecond)
	}
}