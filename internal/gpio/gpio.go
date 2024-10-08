// General Purpose I/O (GPIO) library.
package gpio

import (
	"log"
	"slices"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

var GpioPins = []int{4, 5, 6, 12, 13, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27}

func Read(pin int) rpio.State {
	if err := rpio.Open(); err != nil {
		log.Fatal(err)
	}
	defer rpio.Close()

	p := rpio.Pin(pin)
	p.Input()
	p.PullDown()
	return p.Read()
}

func Write(pin int, value int) {
	if !slices.Contains(GpioPins, pin) {
		log.Fatalf("%d is not a valid GPIO pin", pin)
	}

	if err := rpio.Open(); err != nil {
		log.Fatal(err)
	}
	defer rpio.Close()

	p := rpio.Pin(pin)
	p.Output()
	p.Write(toState(value))
}

func On(pin int, duration time.Duration) {
	Write(pin, 1)
	time.Sleep(duration)
	Off(pin)
}

func Off(pin int) {
	Write(pin, 0)
}

func toState(value int) rpio.State {
	if value == 0 {
		return rpio.Low
	}
	return rpio.High
}
