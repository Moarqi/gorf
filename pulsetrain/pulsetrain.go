package pulsetrain

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

const (
	pinNumber   = 23
	pulseLength = time.Microsecond * 300
)

var (
	pin  rpio.Pin
	open bool = false
)

type pulse struct {
	length int
}

type pulseWagon struct {
	first  pulse
	second pulse
	third  pulse
	fourth pulse
}

// Init initialize gpio access and set desired pin to output mode
func Init() {
	pin = rpio.Pin(pinNumber)
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pin.Output()
	pin.Low()
	open = true
}

// Send sends the decimal code in the pulsetrain format
func Send(id, unit int, state, all bool) bool {
	if !open {
		fmt.Println("error::Send::gpio::not_open")
		os.Exit(1)
	}
	// close gpio access when exiting Send
	// defer rpio.Close()

	var builder strings.Builder
	builder.WriteString(strconv.FormatInt(int64(id), 2))
	if all {
		builder.WriteString(strconv.FormatInt(1, 2))
	} else {
		builder.WriteString(strconv.FormatInt(0, 2))
	}
	if state {
		builder.WriteString(strconv.FormatInt(1, 2))
	} else {
		builder.WriteString(strconv.FormatInt(0, 2))
	}
	builder.WriteString(strconv.FormatInt(int64(unit), 2))

	getRune := func(r rune) rune { return r }
	binaryCodeSlices := strings.Map(getRune, builder.String())
	fmt.Println(binaryCodeSlices)
	// TODO: catch errors when sending
	sendInitBit()
	for i := 0; i < 5; i++ {
		for _, r := range binaryCodeSlices {
			switch r {
			case '0':
				sendOffBit()
			case '1':
				sendOnBit()
			}
		}
		sendFinalBit()
	}

	return true
}

// sendOnBit sends four signals where the second one is high
func sendOnBit() bool {
	sendWaveform(pulseLength, 4*pulseLength)
	sendWaveform(pulseLength, pulseLength/2)
	return true
}

// sendOffBit sends four signals where the last one is high
func sendOffBit() bool {
	sendWaveform(pulseLength, pulseLength/2)
	sendWaveform(pulseLength, 4*pulseLength)
	return true
}

// sendInitBit sends the initializing bit that will be followed by the 32 bit message
func sendInitBit() bool {
	sendWaveform(pulseLength, 8*pulseLength)
	return true
}

// sendFinalBit sends the trailing bit after the 32 bit message
func sendFinalBit() bool {
	sendWaveform(pulseLength, 34*pulseLength)
	return true
}

// https://github.com/martinohmann/rfoutlet
func sleepFor(duration time.Duration) {
	now := time.Now()

	for {
		if time.Since(now) >= duration {
			break
		}
	}
}

// sendWaveform sends on for onDuration followed by off for offDuration, then returns
func sendWaveform(onDuration, offDuration time.Duration) bool {
	pin.High()
	sleepFor(onDuration)
	pin.Low()
	sleepFor(offDuration)
	return true
}
