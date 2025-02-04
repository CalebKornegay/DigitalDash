package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"sync"
	"time"

	"github.com/rzetterberg/elmobd"

	"tinygo.org/x/bluetooth"
)

type DigitalDash struct {
	device               *elmobd.Device
	wg                   *sync.WaitGroup
	bt_adapter           *bluetooth.Adapter
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

/*
* Helper functions
 */

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

func log(format string, a ...any) {
	a = append([]any{time.Now().Format("2006-01-02 15:04:05")}, a...)
	format = "(LOG the time is %s): " + format
	fmt.Printf(format, a...)
}

func Float32ToByte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func Float32FromBytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func connect_to_car() (*elmobd.Device, error) {
	var device *elmobd.Device
	var err error

	device, err = elmobd.NewDevice("/dev/ttyUSB0", false)
	if err == nil {
		return device, nil
	}

	device, err = elmobd.NewDevice("/dev/ttyUSB1", false)
	if err == nil {
		return device, nil
	}

	device, err = elmobd.NewDevice("/dev/ttyUSB2", false)
	if err == nil {
		return device, nil
	}

	device, err = elmobd.NewDevice("/dev/ttyUSB3", false)
	if err == nil {
		return device, nil
	}

	return nil, fmt.Errorf("NO")

}

/*
* Update the values functions
 */

func (dash *DigitalDash) updateRPM(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()

	for {
		cmd := elmobd.NewEngineRPM()
		_, err, _ := dash.device.RunOBDCommand(cmd)

		check_err(err)
		// log("The engine rpm is %f\n", cmd.Value)

		_, err = measurement.Write(Float32ToByte(cmd.Value))
		check_err(err)

		time.Sleep(dash.rpm_wait)
	}
}

func (dash *DigitalDash) updateFuelLevel(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewFuel()
		_, err, _ := dash.device.RunOBDCommand(cmd)

		check_err(err)
		// log("The fuel level is %f%%\n", cmd.Value*100)

		_, err = measurement.Write(Float32ToByte(cmd.Value * 100))
		check_err(err)

		time.Sleep(dash.fuel_level_wait)
	}
}

func (dash *DigitalDash) updateCoolantTemp(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewCoolantTemperature()
		_, err, _ := dash.device.RunOBDCommand(cmd)

		check_err(err)
		// log("The coolant temp is %d\u00b0\n", cmd.Value)

		_, err = measurement.Write(Float32ToByte(float32(cmd.Value)))
		check_err(err)

		time.Sleep(dash.coolant_temp_wait)
	}
}

func (dash *DigitalDash) updateEngineOilTemp(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewEngineOilTemperature()
		_, err, raw_val := dash.device.RunOBDCommand(cmd)

		check_err(err)
		log("The engine oil temp is %d\u00b0\n", cmd.Value)
		log("The raw engine oil temp is %v\n", raw_val)

		_, err = measurement.Write(Float32ToByte(float32(cmd.Value)))
		check_err(err)

		time.Sleep(dash.engine_oil_temp_wait)
	}
}

func (dash *DigitalDash) updateIntakeAirTemp(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewIntakeAirTemperature()
		_, err, _ := dash.device.RunOBDCommand(cmd)

		check_err(err)
		// log("The intake air temp is %d\u00b0C\n", cmd.Value)

		_, err = measurement.Write(Float32ToByte(float32(cmd.Value)))
		check_err(err)

		time.Sleep(dash.intake_air_temp_wait)
	}
}

func (dash *DigitalDash) updateMAFFlowRate(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewMafAirFlowRate()
		_, err, _ := dash.device.RunOBDCommand(cmd)

		check_err(err)
		// log("The mass air flow sensor air flow rate is is %fg/min\n", cmd.Value)

		_, err = measurement.Write(Float32ToByte(cmd.Value))
		check_err(err)

		time.Sleep(dash.maf_wait)
	}
}

func (dash *DigitalDash) updateActualGear(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewTransmissionActualGear()
		_, err, raw_val := dash.device.RunOBDCommand(cmd)

		check_err(err)
		log("The current gear is %s\n", cmd.ValueAsLit())
		log("The raw current gear is %v\n", raw_val)

		_, err = measurement.Write(raw_val)
		check_err(err)

		time.Sleep(dash.gear_wait)
	}
}

func (dash *DigitalDash) updateSpeed(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewVehicleSpeed()
		_, err, _ := dash.device.RunOBDCommand(cmd)

		check_err(err)
		// log("The current speed is %d\n", cmd.Value)

		_, err = measurement.Write(Float32ToByte(float32(cmd.Value & 255)))
		check_err(err)

		time.Sleep(dash.speed_wait)
	}
}

func (dash *DigitalDash) updateAmbientTemp(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewAmbientTemperature()
		_, err, _ := dash.device.RunOBDCommand(cmd)

		check_err(err)
		// log("The ambient temp is %d\u00b0\n", cmd.Value)

		_, err = measurement.Write(Float32ToByte(float32(cmd.Value)))
		check_err(err)

		time.Sleep(dash.ambient_temp_wait)
	}
}

func (dash *DigitalDash) updateThrottlePosition(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewThrottlePosition()
		_, err, _ := dash.device.RunOBDCommand(cmd)

		check_err(err)
		// log("The throttle position is %f%%\n", cmd.Value*100)

		_, err = measurement.Write(Float32ToByte(cmd.Value * 100))
		check_err(err)

		time.Sleep(dash.throttle_pos_wait)
	}
}

func (dash *DigitalDash) updateOdometer(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewOdometer()
		_, err, _ := dash.device.RunOBDCommand(cmd)

		check_err(err)
		// log("The current mileage is %f miles\n", cmd.Value)

		_, err = measurement.Write(Float32ToByte(cmd.Value))
		check_err(err)

		time.Sleep(dash.odometer_wait)
	}
}

func (dash *DigitalDash) updateVoltage(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()

	for {
		voltage, err := dash.device.GetVoltage()

		check_err(err)

		// log("The current battery voltage is %fV\n", voltage)

		_, err = measurement.Write(Float32ToByte(voltage))
		check_err(err)

		time.Sleep(dash.voltage_wait)
	}
}

func main() {
	var device *elmobd.Device
	var err error

	adapter := bluetooth.DefaultAdapter
	err = adapter.Enable()
	fatal(err)

	var rpmMeasurement, coolant_tempMeasurement, intake_air_tempMeasurement, speedMeasurement, ambient_tempMeasurement, fuel_levelMeasurement, maf_flow_rateMeasurement, throttle_posMeasurement, voltageMeasurement, engine_oil_tempMeasurement, gearMeasurement, odometerMeasurement bluetooth.Characteristic

	// 0x272F degrees celsius
	// 0x2767 volume liters
	// 0x27A4 distance miles
	// 0x27A7 velocity mph
	// 0x27AD percentage
	// 0x27AF period revs per minute
	// 0x27C1 mass flow grams per second
	// 0x27C2 volume flow liters per second

	advertisement := adapter.DefaultAdvertisement()
	err = advertisement.Configure(bluetooth.AdvertisementOptions{
		LocalName:    "Raspberry Pi OBD-II",
		ServiceUUIDs: []bluetooth.UUID{bluetooth.ServiceUUIDHumanInterfaceDevice},
	})
	fatal(err)

	err = advertisement.Start()
	fatal(err)

	err = adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDHumanInterfaceDevice,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &rpmMeasurement,
				UUID:   bluetooth.New16BitUUID(0x27AF), // revs/min
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
			},
			{
				Handle: &coolant_tempMeasurement,
				UUID:   bluetooth.New16BitUUID(0x272F), // degrees C
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
			},
			{
				Handle: &intake_air_tempMeasurement,
				UUID:   bluetooth.New16BitUUID(0x2730), // degrees C
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
			},
			{
				Handle: &speedMeasurement,
				UUID:   bluetooth.New16BitUUID(0x27A7), // speed mph
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
			},
			{
				Handle: &ambient_tempMeasurement,
				UUID:   bluetooth.New16BitUUID(0x2731), // degrees C
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
			},
			{
				Handle: &fuel_levelMeasurement,
				UUID:   bluetooth.New16BitUUID(0x27AD), // percentage
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
			},
			{
				Handle: &maf_flow_rateMeasurement,
				UUID:   bluetooth.New16BitUUID(0x27C1), // flow grams/sec
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
			},
			{
				Handle: &throttle_posMeasurement,
				UUID:   bluetooth.New16BitUUID(0x27AE), // percentage
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
			},
			{
				Handle: &voltageMeasurement,
				UUID:   bluetooth.New16BitUUID(0x2B18), // voltage
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
			},
			{
				Handle: &odometerMeasurement,
				UUID:   bluetooth.New16BitUUID(0x27A4), // distance miles
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
			},
			{
				Handle: &engine_oil_tempMeasurement,
				UUID:   bluetooth.New16BitUUID(0x2732), // degrees C
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
			},
			{
				Handle: &gearMeasurement,
				UUID:   bluetooth.New16BitUUID(0x2C0B), // torque (gear ratio)
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
			},
		},
	})
	fatal(err)

	// var chars = []bluetooth.CharacteristicConfig{
	// 	{
	// 		Handle: &rpmMeasurement,
	// 		UUID:   bluetooth.New16BitUUID(0x27AF), // revs/min
	// 		Value:  []byte{0, 0, 0, 0},
	// 		Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
	// 	},
	// 	{
	// 		Handle: &coolant_tempMeasurement,
	// 		UUID:   bluetooth.New16BitUUID(0x272F), // degrees C
	// 		Value:  []byte{0, 0, 0, 0},
	// 		Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
	// 	},
	// 	{
	// 		Handle: &intake_air_tempMeasurement,
	// 		UUID:   bluetooth.New16BitUUID(0x2730), // degrees C
	// 		Value:  []byte{0, 0, 0, 0},
	// 		Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
	// 	},
	// 	{
	// 		Handle: &speedMeasurement,
	// 		UUID:   bluetooth.New16BitUUID(0x27A7), // speed mph
	// 		Value:  []byte{0, 0, 0, 0},
	// 		Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
	// 	},
	// 	{
	// 		Handle: &ambient_tempMeasurement,
	// 		UUID:   bluetooth.New16BitUUID(0x2731), // degrees C
	// 		Value:  []byte{0, 0, 0, 0},
	// 		Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
	// 	},
	// 	{
	// 		Handle: &fuel_levelMeasurement,
	// 		UUID:   bluetooth.New16BitUUID(0x27AD), // percentage
	// 		Value:  []byte{0, 0, 0, 0},
	// 		Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
	// 	},
	// 	{
	// 		Handle: &maf_flow_rateMeasurement,
	// 		UUID:   bluetooth.New16BitUUID(0x27C1), // flow grams/sec
	// 		Value:  []byte{0, 0, 0, 0},
	// 		Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
	// 	},
	// 	{
	// 		Handle: &throttle_posMeasurement,
	// 		UUID:   bluetooth.New16BitUUID(0x27AD), // percentage
	// 		Value:  []byte{0, 0, 0, 0},
	// 		Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
	// 	},
	// 	{
	// 		Handle: &voltageMeasurement,
	// 		UUID:   bluetooth.New16BitUUID(0x2B18), // voltage
	// 		Value:  []byte{0, 0, 0, 0},
	// 		Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
	// 	},
	// 	{
	// 		Handle: &odometerMeasurement,
	// 		UUID:   bluetooth.New16BitUUID(0x27A4), // distance miles
	// 		Value:  []byte{0, 0, 0, 0},
	// 		Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
	// 	},
	// 	{
	// 		Handle: &engine_oil_tempMeasurement,
	// 		UUID:   bluetooth.New16BitUUID(0x2732), // degrees C
	// 		Value:  []byte{0, 0, 0, 0},
	// 		Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
	// 	},
	// 	{
	// 		Handle: &gearMeasurement,
	// 		UUID:   bluetooth.New16BitUUID(0x2C0B), // torque (gear ratio)
	// 		Value:  []byte{0, 0, 0, 0},
	// 		Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
	// 	},
	// }

	// var i float32 = 0.0
	// for {
	// 	for char := range chars {
	// 		chars[char].Handle.Write(Float32ToByte(i))
	// 	}
	// 	i += 1
	// 	time.Sleep(time.Second * 1)
	// }

	// Try to connect to the device multiple times before giving up
	for {
		device, err = connect_to_car()
		if err == nil {
			break
		}
		time.Sleep(time.Second * 3)
	}

	time.Sleep(time.Second * 3) // Let the device initialize
	wg := sync.WaitGroup{}

	var dash DigitalDash = DigitalDash{
		device:               device,
		wg:                   &wg,
		bt_adapter:           adapter,
		rpm_wait:             time.Millisecond * 100,
		coolant_temp_wait:    time.Millisecond * 500,
		engine_oil_temp_wait: time.Millisecond * 500,
		intake_air_temp_wait: time.Millisecond * 500,
		maf_wait:             time.Millisecond * 500,
		speed_wait:           time.Millisecond * 100,
		throttle_pos_wait:    time.Millisecond * 250,
		fuel_level_wait:      time.Second * 5,
		odometer_wait:        time.Second * 5,
		voltage_wait:         time.Second * 5,
	}

	fmt.Printf("The car is on and the time is %s\n", time.Now())

	wg.Add(12)
	go dash.updateRPM(&rpmMeasurement)
	go dash.updateFuelLevel(&fuel_levelMeasurement)
	go dash.updateCoolantTemp(&coolant_tempMeasurement)
	go dash.updateEngineOilTemp(&engine_oil_tempMeasurement) // ECHO mismatch
	go dash.updateIntakeAirTemp(&intake_air_tempMeasurement)
	go dash.updateMAFFlowRate(&maf_flow_rateMeasurement)
	go dash.updateOdometer(&odometerMeasurement)
	go dash.updateSpeed(&speedMeasurement)
	go dash.updateThrottlePosition(&throttle_posMeasurement)
	go dash.updateActualGear(&gearMeasurement) // NO DATA ??
	go dash.updateAmbientTemp(&ambient_tempMeasurement)
	go dash.updateVoltage(&voltageMeasurement)

	wg.Wait()
}
