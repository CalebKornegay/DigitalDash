package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/rzetterberg/elmobd"
)

type DigitalDash struct {
	rpm_wait        time.Duration
	fuel_level_wait time.Duration
	device          *elmobd.Device
	wg              sync.WaitGroup
}

func try_panic(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error was %s... Exiting\n", err.Error())
		os.Exit(1)
	}
}

func (dash *DigitalDash) updateRPM() {
	dash.wg.Add(1)

	for {
		cmd := elmobd.NewEngineRPM()
		_, err := dash.device.RunOBDCommand(cmd)

		try_panic(err)
		fmt.Printf("The engine rpm is %f\n", cmd.Value)

		time.Sleep(dash.rpm_wait)
	}

	dash.wg.Done()
}

func (dash *DigitalDash) updateFuelLevel() {
	dash.wg.Add(1)

	for {
		cmd := elmobd.NewFuel()
		_, err := dash.device.RunOBDCommand(cmd)

		try_panic(err)
		fmt.Printf("The fuel level is %f%%\n", cmd.Value)
	}

	dash.wg.Done()
}

func main() {

	device, err := elmobd.NewDevice("/dev/ttyUSB0", true)
	try_panic(err)

	var dash DigitalDash = DigitalDash{
		device:          device,
		rpm_wait:        time.Millisecond * 50,
		fuel_level_wait: time.Second * 5,
	}

	// on, err := device.GetIgnitionState()
	// try_panic(err)
	// fmt.Printf("The car is on: %t\n", on)

	go dash.updateRPM()
	go dash.updateFuelLevel()

	dash.wg.Wait()
}
