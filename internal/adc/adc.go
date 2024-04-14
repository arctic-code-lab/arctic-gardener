// Analog to Digital Converter (ADC) library.
package adc

import (
	"log"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/ads1x15"
	"periph.io/x/host/v3"
)

const (
	AdcMaxValue = 26400 // ADS1115 gives you 32768 steps with 4.096v, with 3.3v which is the 80% of 4.096v we get 26400 steps
	AdcMinValue = 0
	inVoltage   = 3300 * physic.MilliVolt
	inFrequency = 1 * physic.Hertz
)

func Read(channel int) int {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Open default I²C bus.
	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatalf("failed to open I²C: %v", err)
	}
	defer bus.Close()

	// Create a new ADS1115 ADC.
	adc, err := ads1x15.NewADS1115(bus, &ads1x15.DefaultOpts)
	if err != nil {
		log.Fatalln(err)
	}

	// Obtain an analog pin from the ADC.
	pin, err := adc.PinForChannel(getChannel(channel), inVoltage, inFrequency, ads1x15.SaveEnergy)
	if err != nil {
		log.Fatalln(err)
	}
	defer pin.Halt()

	// Read values from ADC.
	reading, err := pin.Read()
	if err != nil {
		log.Fatalln(err)
	}

	return int(reading.Raw)
}

func getChannel(channel int) ads1x15.Channel {
	switch channel {
	case 0:
		return ads1x15.Channel0
	case 1:
		return ads1x15.Channel1
	case 2:
		return ads1x15.Channel2
	case 3:
		return ads1x15.Channel3
	default:
		return ads1x15.Channel0
	}
}
