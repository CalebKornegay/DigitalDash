package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/rzetterberg/elmobd"
)

type DigitalDash struct {
	device               *elmobd.Device
	wg                   sync.WaitGroup
	rpm_wait             time.Duration
	fuel_level_wait      time.Duration
	coolant_temp_wait    time.Duration
	engine_oil_temp_wait time.Duration
	intake_air_temp_wait time.Duration
	maf_wait             time.Duration
	gear_wait            time.Duration
	speed_wait           time.Duration
	ambient_temp_wait    time.Duration
	throttle_pos_wait    time.Duration
	odometer_wait        time.Duration
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

func (dash *DigitalDash) updateCoolantTemp() {
	dash.wg.Add(1)

	for {
		cmd := elmobd.NewCoolantTemperature()
		_, err := dash.device.RunOBDCommand(cmd)

		try_panic(err)
		fmt.Printf("The coolant temp is %d\u00b0\n", cmd.Value)
	}

	dash.wg.Done()
}

func (dash *DigitalDash) updateEngineOilTemp() {
	dash.wg.Add(1)

	for {
		cmd := elmobd.NewEngineOilTemperature()
		_, err := dash.device.RunOBDCommand(cmd)

		try_panic(err)
		fmt.Printf("The engine oil temp is %d\u00b0\n", cmd.Value)
	}

	dash.wg.Done()
}

func (dash *DigitalDash) updateIntakeAirTemp() {
	dash.wg.Add(1)

	for {
		cmd := elmobd.NewIntakeAirTemperature()
		_, err := dash.device.RunOBDCommand(cmd)

		try_panic(err)
		fmt.Printf("The intake air temp is %d\u00b0\n", cmd.Value)
	}

	dash.wg.Done()
}

func (dash *DigitalDash) updateMAFFlowRate() {
	dash.wg.Add(1)

	for {
		cmd := elmobd.NewMafAirFlowRate()
		_, err := dash.device.RunOBDCommand(cmd)

		try_panic(err)
		fmt.Printf("The mass air flow sensor air flow rate is is %fL/hr\n", cmd.Value)
	}

	dash.wg.Done()
}

func (dash *DigitalDash) updateActualGear() {
	dash.wg.Add(1)

	for {
		cmd := elmobd.NewTransmissionActualGear()
		_, err := dash.device.RunOBDCommand(cmd)

		try_panic(err)
		fmt.Printf("The current gear is %s\n", cmd.ValueAsLit())
	}

	dash.wg.Done()
}

func (dash *DigitalDash) updateSpeed() {
	dash.wg.Add(1)

	for {
		cmd := elmobd.NewVehicleSpeed()
		_, err := dash.device.RunOBDCommand(cmd)

		try_panic(err)
		fmt.Printf("The current speed is %d\n", cmd.Value)
	}

	dash.wg.Done()
}

func (dash *DigitalDash) updateAmbientTemp() {
	dash.wg.Add(1)

	for {
		cmd := elmobd.NewAmbientTemperature()
		_, err := dash.device.RunOBDCommand(cmd)

		try_panic(err)
		fmt.Printf("The ambient temp is %d\u00b0\n", cmd.Value)
	}

	dash.wg.Done()
}

func (dash *DigitalDash) updateThrottlePosition() {
	dash.wg.Add(1)

	for {
		cmd := elmobd.NewThrottlePosition()
		_, err := dash.device.RunOBDCommand(cmd)

		try_panic(err)
		fmt.Printf("The throttle position is %f%%\n", cmd.Value)
	}

	dash.wg.Done()
}

func (dash *DigitalDash) updateOdometer() {
	dash.wg.Add(1)

	for {
		cmd := elmobd.NewOdometer()
		_, err := dash.device.RunOBDCommand(cmd)

		try_panic(err)
		fmt.Printf("The current mileage is %f miles\n", cmd.Value)
	}

	dash.wg.Done()
}

func main() {

	device, err := elmobd.NewDevice("/dev/ttyUSB0", true)
	try_panic(err)

	var dash DigitalDash = DigitalDash{
		device:               device,
		rpm_wait:             time.Millisecond * 50,
		fuel_level_wait:      time.Second * 5,
		coolant_temp_wait:    time.Millisecond * 100,
		engine_oil_temp_wait: time.Millisecond * 100,
		intake_air_temp_wait: time.Millisecond * 100,
		maf_wait:             time.Millisecond * 500,
		gear_wait:            time.Second * 1,
		speed_wait:           time.Millisecond * 50,
		ambient_temp_wait:    time.Millisecond * 500,
		throttle_pos_wait:    time.Millisecond * 500,
		odometer_wait:        time.Second * 5,
	}

	// on, err := device.GetIgnitionState()
	// try_panic(err)
	// fmt.Printf("The car is on: %t\n", on)

	go dash.updateRPM()
	go dash.updateFuelLevel()
	go dash.updateCoolantTemp()
	go dash.updateEngineOilTemp()
	go dash.updateIntakeAirTemp()
	go dash.updateMAFFlowRate()
	go dash.updateOdometer()
	go dash.updateSpeed()
	go dash.updateThrottlePosition()
	go dash.updateActualGear()
	go dash.updateAmbientTemp()

	dash.wg.Wait()
}
