package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rzetterberg/elmobd"
)

func try_panic(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error was %s... Exiting\n", err.Error())
		os.Exit(1)
	}
}

func main() {
	time.Sleep(time.Second * 1)

	info, err := os.Stat("/dev/ttyUSB0")
	try_panic(err)
	fmt.Println(info)

	device, err := elmobd.NewDevice("/dev/ttyUSB0", true)
	try_panic(err)

	on, err := device.GetIgnitionState()
	try_panic(err)
	fmt.Printf("The car is on: %t\n", on)

	rpm, err := device.RunOBDCommand(elmobd.NewEngineRPM())
	try_panic(err)
	fmt.Printf("The engine rpm is %s\n", rpm.ValueAsLit())

	fuel, err := device.RunOBDCommand(elmobd.NewFuel())
	try_panic(err)
	fmt.Printf("The fuel level is %s\n", fuel.ValueAsLit())

}
