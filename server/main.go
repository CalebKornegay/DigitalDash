package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

func Float32ToByte(f float32) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, f)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
	// tmp := uint32(f)
	// return []byte{uint8(tmp >> 24), uint8(tmp >> 16), uint8(tmp >> 8), uint8(tmp)}
}

func (dash *DigitalDash) updateRPM(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()

	for {
		cmd := elmobd.NewEngineRPM()
		_, err := dash.device.RunOBDCommand(cmd)

		check_err(err)
		log("The engine rpm is %f\n", cmd.Value)

		_, err = measurement.Write(Float32ToByte(cmd.Value))
		check_err(err)

		time.Sleep(dash.rpm_wait)
	}
}

func (dash *DigitalDash) updateFuelLevel(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewFuel()
		_, err := dash.device.RunOBDCommand(cmd)

		check_err(err)
		log("The fuel level is %f%%\n", cmd.Value*100)

		_, err = measurement.Write(Float32ToByte(cmd.Value * 100))
		check_err(err)

		time.Sleep(dash.fuel_level_wait)
	}
}

func (dash *DigitalDash) updateCoolantTemp(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewCoolantTemperature()
		_, err := dash.device.RunOBDCommand(cmd)

		check_err(err)
		log("The coolant temp is %d\u00b0\n", cmd.Value)

		_, err = measurement.Write(Float32ToByte(float32(cmd.Value)))
		check_err(err)

		time.Sleep(dash.coolant_temp_wait)
	}
}

func (dash *DigitalDash) updateEngineOilTemp(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewEngineOilTemperature()
		_, err := dash.device.RunOBDCommand(cmd)

		check_err(err)
		log("The engine oil temp is %d\u00b0\n", cmd.Value)

		_, err = measurement.Write(Float32ToByte(float32(cmd.Value)))
		check_err(err)

		time.Sleep(dash.engine_oil_temp_wait)
	}
}

func (dash *DigitalDash) updateIntakeAirTemp(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewIntakeAirTemperature()
		_, err := dash.device.RunOBDCommand(cmd)

		check_err(err)
		log("The intake air temp is %d\u00b0C\n", cmd.Value)

		_, err = measurement.Write(Float32ToByte(float32(cmd.Value)))
		check_err(err)

		time.Sleep(dash.intake_air_temp_wait)
	}
}

func (dash *DigitalDash) updateMAFFlowRate(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewMafAirFlowRate()
		_, err := dash.device.RunOBDCommand(cmd)

		check_err(err)
		log("The mass air flow sensor air flow rate is is %fg/min\n", cmd.Value)

		_, err = measurement.Write(Float32ToByte(cmd.Value))
		check_err(err)

		time.Sleep(dash.maf_wait)
	}
}

func (dash *DigitalDash) updateActualGear(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewTransmissionActualGear()
		_, err := dash.device.RunOBDCommand(cmd)

		check_err(err)
		log("The current gear is %s\n", cmd.ValueAsLit())

		_, err = measurement.Write(Float32ToByte(cmd.Value))
		check_err(err)

		time.Sleep(dash.gear_wait)
	}
}

func (dash *DigitalDash) updateSpeed(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewVehicleSpeed()
		_, err := dash.device.RunOBDCommand(cmd)

		check_err(err)
		log("The current speed is %d\n", cmd.Value)

		_, err = measurement.Write(Float32ToByte(float32(cmd.Value)))
		check_err(err)

		time.Sleep(dash.speed_wait)
	}
}

func (dash *DigitalDash) updateAmbientTemp(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewAmbientTemperature()
		_, err := dash.device.RunOBDCommand(cmd)

		check_err(err)
		log("The ambient temp is %d\u00b0\n", cmd.Value)

		_, err = measurement.Write(Float32ToByte(float32(cmd.Value)))
		check_err(err)

		time.Sleep(dash.ambient_temp_wait)
	}
}

func (dash *DigitalDash) updateThrottlePosition(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewThrottlePosition()
		_, err := dash.device.RunOBDCommand(cmd)

		check_err(err)
		log("The throttle position is %f%%\n", cmd.Value*100)

		_, err = measurement.Write(Float32ToByte(cmd.Value * 100))
		check_err(err)

		time.Sleep(dash.throttle_pos_wait)
	}
}

func (dash *DigitalDash) updateOdometer(measurement *bluetooth.Characteristic) {
	defer dash.wg.Done()
	for {
		cmd := elmobd.NewOdometer()
		_, err := dash.device.RunOBDCommand(cmd)

		check_err(err)
		log("The current mileage is %f miles\n", cmd.Value)

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

		log("The current battery voltage is %fV\n", voltage)

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

	var rpmMeasurement, coolant_tempMeasurement, intake_air_tempMeasurement, speedMeasurement, ambient_tempMeasurement, fuel_levelMeasurement, maf_flow_rateMeasurement, throttle_posMeasurement, voltageMeasurement, engine_oil_tempMeasuremnt, gearMeasurement, odometerMeasurement bluetooth.Characteristic

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
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
			},
			{
				Handle: &coolant_tempMeasurement,
				UUID:   bluetooth.New16BitUUID(0x272F), // degrees C
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
			},
			{
				Handle: &intake_air_tempMeasurement,
				UUID:   bluetooth.New16BitUUID(0x2730), // degrees C
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
			},
			{
				Handle: &speedMeasurement,
				UUID:   bluetooth.New16BitUUID(0x27A7), // speed mph
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
			},
			{
				Handle: &ambient_tempMeasurement,
				UUID:   bluetooth.New16BitUUID(0x2731), // degrees C
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
			},
			{
				Handle: &fuel_levelMeasurement,
				UUID:   bluetooth.New16BitUUID(0x27AD), // percentage
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
			},
			{
				Handle: &maf_flow_rateMeasurement,
				UUID:   bluetooth.New16BitUUID(0x27C1), // flow grams/sec
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
			},
			{
				Handle: &throttle_posMeasurement,
				UUID:   bluetooth.New16BitUUID(0x27AD), // percentage
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
			},
			{
				Handle: &voltageMeasurement,
				UUID:   bluetooth.New16BitUUID(0x2B18), // voltage
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
			},
			{
				Handle: &odometerMeasurement,
				UUID:   bluetooth.New16BitUUID(0x27A4), // distance miles
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
			},
			{
				Handle: &engine_oil_tempMeasuremnt,
				UUID:   bluetooth.New16BitUUID(0x2732), // degrees C
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
			},
			{
				Handle: &gearMeasurement,
				UUID:   bluetooth.New16BitUUID(0x2C0B), // torque (gear ratio)
				Value:  []byte{0, 0, 0, 0},
				Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
			},
		},
	})
	fatal(err)

	var chars = []bluetooth.CharacteristicConfig{
		{
			Handle: &rpmMeasurement,
			UUID:   bluetooth.New16BitUUID(0x27AF), // revs/min
			Value:  []byte{0, 0, 0, 0},
			Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
		},
		{
			Handle: &coolant_tempMeasurement,
			UUID:   bluetooth.New16BitUUID(0x272F), // degrees C
			Value:  []byte{0, 0, 0, 0},
			Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
		},
		{
			Handle: &intake_air_tempMeasurement,
			UUID:   bluetooth.New16BitUUID(0x2730), // degrees C
			Value:  []byte{0, 0, 0, 0},
			Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
		},
		{
			Handle: &speedMeasurement,
			UUID:   bluetooth.New16BitUUID(0x27A7), // speed mph
			Value:  []byte{0, 0, 0, 0},
			Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
		},
		{
			Handle: &ambient_tempMeasurement,
			UUID:   bluetooth.New16BitUUID(0x2731), // degrees C
			Value:  []byte{0, 0, 0, 0},
			Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
		},
		{
			Handle: &fuel_levelMeasurement,
			UUID:   bluetooth.New16BitUUID(0x27AD), // percentage
			Value:  []byte{0, 0, 0, 0},
			Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
		},
		{
			Handle: &maf_flow_rateMeasurement,
			UUID:   bluetooth.New16BitUUID(0x27C1), // flow grams/sec
			Value:  []byte{0, 0, 0, 0},
			Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
		},
		{
			Handle: &throttle_posMeasurement,
			UUID:   bluetooth.New16BitUUID(0x27AD), // percentage
			Value:  []byte{0, 0, 0, 0},
			Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
		},
		{
			Handle: &voltageMeasurement,
			UUID:   bluetooth.New16BitUUID(0x2B18), // voltage
			Value:  []byte{0, 0, 0, 0},
			Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
		},
		{
			Handle: &odometerMeasurement,
			UUID:   bluetooth.New16BitUUID(0x27A4), // distance miles
			Value:  []byte{0, 0, 0, 0},
			Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
		},
		{
			Handle: &engine_oil_tempMeasuremnt,
			UUID:   bluetooth.New16BitUUID(0x2732), // degrees C
			Value:  []byte{0, 0, 0, 0},
			Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
		},
		{
			Handle: &gearMeasurement,
			UUID:   bluetooth.New16BitUUID(0x2C0B), // torque (gear ratio)
			Value:  []byte{0, 0, 0, 0},
			Flags:  bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicNotifyPermission,
		},
	}

	var i byte = 0
	for {
		for char := range chars {
			chars[char].Handle.Write([]byte{i, i, i, i})
			i += 1
		}
	}

	// Try to connect to the device multiple times before giving up
	// for i := 0; i < 5; i++ {
	// 	device, err = elmobd.NewDevice("/dev/ttyUSB0", false)
	// 	if err == nil {
	// 		break
	// 	} else if i == 4 {
	// 		fatal(err)
	// 	}
	// 	time.Sleep(time.Second * 3)
	// }

	// time.Sleep(time.Second * 3) // Let the device initialize

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
		odometer_wait:        time.Second * 5,
		voltage_wait:         time.Second * 5,
	}

	fmt.Printf("The car is on and the time is %s\n", time.Now())

	wg.Add(12)
	go dash.updateRPM(&rpmMeasurement)
	go dash.updateFuelLevel(&fuel_levelMeasurement)
	go dash.updateCoolantTemp(&coolant_tempMeasurement)
	go dash.updateEngineOilTemp(&engine_oil_tempMeasuremnt) // ECHO mismatch
	go dash.updateIntakeAirTemp(&intake_air_tempMeasurement)
	go dash.updateMAFFlowRate(&maf_flow_rateMeasurement)
	go dash.updateOdometer(&odometerMeasurement) // Reads incorrectly (18830.1 miles when I have 11700.2)
	go dash.updateSpeed(&speedMeasurement)
	go dash.updateThrottlePosition(&throttle_posMeasurement)
	go dash.updateActualGear(&gearMeasurement) // NO DATA ??
	go dash.updateAmbientTemp(&ambient_tempMeasurement)
	go dash.updateVoltage(&voltageMeasurement)

	wg.Wait()
}
