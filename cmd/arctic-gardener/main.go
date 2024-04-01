/*
Rev 2 and 3 Raspberry Pi
+-----+---------+----------+---------+-----+
| BCM |   Name  | Physical | Name    | BCM |
+-----+---------+----++----+---------+-----+
|     |    3.3v |  1 || 2  | 5v      |     |
|   2 |   SDA 1 |  3 || 4  | 5v      |     |
|   3 |   SCL 1 |  5 || 6  | 0v      |     |
|   4 | GPIO  7 |  7 || 8  | TxD     | 14  |
|     |      0v |  9 || 10 | RxD     | 15  |
|  17 | GPIO  0 | 11 || 12 | GPIO  1 | 18  |
|  27 | GPIO  2 | 13 || 14 | 0v      |     |
|  22 | GPIO  3 | 15 || 16 | GPIO  4 | 23  |
|     |    3.3v | 17 || 18 | GPIO  5 | 24  |
|  10 |    MOSI | 19 || 20 | 0v      |     |
|   9 |    MISO | 21 || 22 | GPIO  6 | 25  |
|  11 |    SCLK | 23 || 24 | CE0     | 8   |
|     |      0v | 25 || 26 | CE1     | 7   |
|   0 |   SDA 0 | 27 || 28 | SCL 0   | 1   |
|   5 | GPIO 21 | 29 || 30 | 0v      |     |
|   6 | GPIO 22 | 31 || 32 | GPIO 26 | 12  |
|  13 | GPIO 23 | 33 || 34 | 0v      |     |
|  19 | GPIO 24 | 35 || 36 | GPIO 27 | 16  |
|  26 | GPIO 25 | 37 || 38 | GPIO 28 | 20  |
|     |      0v | 39 || 40 | GPIO 29 | 21  |
+-----+---------+----++----+---------+-----+
*/

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

func main() {
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

	pin := rpio.Pin(4)
	pin.Input()
	pin.PullDown()

	for {
		fmt.Println("PullDown:", pin.Read())
		time.Sleep(1500 * time.Millisecond)
	}
}