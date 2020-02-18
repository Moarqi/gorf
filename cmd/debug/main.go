package main

import (
	"fmt"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

const (
	pinNumber = 24
)

func main() {
	var pin rpio.Pin
	pin = rpio.Pin(pinNumber)
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pin.Output()

	for true {
		pin.High()
		time.Sleep(2 * time.Second)
		pin.Low()
		time.Sleep(2 * time.Second)
	}
}
