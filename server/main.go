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
	wg                   *sync.WaitGroup
	lock                 *sync.Mutex
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
	voltage_wait         time.Duration
}

func check_err(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error was %s\n", err.Error())
	}
}

func fatal(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error was %s\n", err.Error())
		os.Exit(1)
	}
}

func (dash *DigitalDash) updateRPM() {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewEngineRPM()
		_, err := dash.device.RunOBDCommand(cmd)

		check_err(err)
		fmt.Printf("The engine rpm is %f\n", cmd.Value)

		time.Sleep(dash.rpm_wait)
	}

}

func (dash *DigitalDash) updateFuelLevel() {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewFuel()
		_, err := dash.device.RunOBDCommand(cmd)

		check_err(err)
		fmt.Printf("The fuel level is %f%%\n", cmd.Value*100)

		time.Sleep(dash.fuel_level_wait)
	}
}

func (dash *DigitalDash) updateCoolantTemp() {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewCoolantTemperature()
		_, err := dash.device.RunOBDCommand(cmd)

		check_err(err)
		fmt.Printf("The coolant temp is %d\u00b0\n", cmd.Value)

		time.Sleep(dash.coolant_temp_wait)
	}
}

func (dash *DigitalDash) updateEngineOilTemp() {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewEngineOilTemperature()
		_, err := dash.device.RunOBDCommand(cmd)

		check_err(err)
		fmt.Printf("The engine oil temp is %d\u00b0\n", cmd.Value)

		time.Sleep(dash.engine_oil_temp_wait)
	}
}

func (dash *DigitalDash) updateIntakeAirTemp() {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewIntakeAirTemperature()
		_, err := dash.device.RunOBDCommand(cmd)

		check_err(err)
		fmt.Printf("The intake air temp is %d\u00b0\n", cmd.Value)

		time.Sleep(dash.intake_air_temp_wait)
	}
}

func (dash *DigitalDash) updateMAFFlowRate() {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewMafAirFlowRate()
		_, err := dash.device.RunOBDCommand(cmd)

		check_err(err)
		fmt.Printf("The mass air flow sensor air flow rate is is %fL/hr\n", cmd.Value)

		time.Sleep(dash.maf_wait)
	}
}

func (dash *DigitalDash) updateActualGear() {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewTransmissionActualGear()
		_, err := dash.device.RunOBDCommand(cmd)

		check_err(err)
		fmt.Printf("The current gear is %s\n", cmd.ValueAsLit())

		time.Sleep(dash.gear_wait)
	}
}

func (dash *DigitalDash) updateSpeed() {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewVehicleSpeed()
		_, err := dash.device.RunOBDCommand(cmd)

		check_err(err)
		fmt.Printf("The current speed is %d\n", cmd.Value)

		time.Sleep(dash.speed_wait)
	}
}

func (dash *DigitalDash) updateAmbientTemp() {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewAmbientTemperature()
		_, err := dash.device.RunOBDCommand(cmd)

		check_err(err)
		fmt.Printf("The ambient temp is %d\u00b0\n", cmd.Value)

		time.Sleep(dash.ambient_temp_wait)
	}
}

func (dash *DigitalDash) updateThrottlePosition() {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewThrottlePosition()
		_, err := dash.device.RunOBDCommand(cmd)

		check_err(err)
		fmt.Printf("The throttle position is %f%%\n", cmd.Value*100)

		time.Sleep(dash.throttle_pos_wait)
	}
}

func (dash *DigitalDash) updateOdometer() {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewOdometer()
		_, err := dash.device.RunOBDCommand(cmd)

		check_err(err)
		fmt.Printf("The current mileage is %f miles\n", cmd.Value)

		time.Sleep(dash.odometer_wait)
	}
}

func (dash *DigitalDash) updateVoltage() {
	defer dash.wg.Done()

	for {
		voltage, err := dash.device.GetVoltage()

		check_err(err)

		fmt.Printf("The current battery voltage is %fV\n", voltage)
		time.Sleep(dash.voltage_wait)
	}
}

func main() {

	device, err := elmobd.NewDevice("/dev/ttyUSB0", true)
	fatal(err)

	time.Sleep(time.Second * 2) // Let the device initialize

	// supported, err := device.CheckSupportedCommands()

	// if err != nil {
	// 	fmt.Println("Failed to check supported commands", err)
	// 	return
	// }

	// allCommands := elmobd.GetSensorCommands()
	// carCommands := supported.FilterSupported(allCommands)

	// fmt.Printf("%d of %d commands supported:\n", len(carCommands), len(allCommands))

	// for _, cmd := range carCommands {
	// 	fmt.Printf("- %s supported\n", cmd.Key())
	// }

	wg := sync.WaitGroup{}
	lock := sync.Mutex{}

	var dash DigitalDash = DigitalDash{
		device:               device,
		wg:                   &wg,
		lock:                 &lock,
		rpm_wait:             time.Millisecond * 50,
		fuel_level_wait:      time.Second * 5,
		coolant_temp_wait:    time.Millisecond * 100,
		engine_oil_temp_wait: time.Millisecond * 100,
		intake_air_temp_wait: time.Millisecond * 100,
		maf_wait:             time.Millisecond * 500,
		speed_wait:           time.Millisecond * 50,
		throttle_pos_wait:    time.Millisecond * 500,
		odometer_wait:        time.Second * 5,
		voltage_wait:         time.Second * 5,
	}

	// on, err := device.GetIgnitionState()
	// check_err(err)
	// fmt.Printf("The car is on: %t\n", on)

	wg.Add(12)
	go dash.updateRPM()
	go dash.updateFuelLevel()
	go dash.updateCoolantTemp()
	// go dash.updateEngineOilTemp() // ECHO mismatch
	go dash.updateIntakeAirTemp()
	go dash.updateMAFFlowRate()
	// go dash.updateOdometer() // Reads incorrectly (18830.1 miles when I have 11700.2)
	go dash.updateSpeed()
	go dash.updateThrottlePosition()
	// go dash.updateActualGear() // NO DATA ??
	go dash.updateAmbientTemp()
	go dash.updateVoltage()

	wg.Wait()
}
