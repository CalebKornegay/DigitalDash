package com.example.digitaldash

import android.Manifest
import android.app.ActionBar
import android.bluetooth.BluetoothAdapter
import android.bluetooth.BluetoothGatt
import android.bluetooth.BluetoothGattCallback
import android.bluetooth.BluetoothGattCharacteristic
import android.bluetooth.BluetoothGattDescriptor
import android.bluetooth.BluetoothManager
import android.bluetooth.BluetoothProfile
import android.content.Context
import android.content.pm.PackageManager
import android.os.Bundle
import android.util.Log
import android.view.View
import android.view.Window
import android.widget.Button
import android.widget.TextView
import androidx.activity.ComponentActivity
import androidx.activity.enableEdgeToEdge
import androidx.core.app.ActivityCompat
import androidx.core.content.ContextCompat
import java.util.Locale
import java.util.Queue
import java.util.UUID
import java.util.concurrent.ConcurrentLinkedQueue


class MainActivity : ComponentActivity() {
    private val bluetoothAdapter: BluetoothAdapter by lazy {
        val bluetoothManager = getSystemService(Context.BLUETOOTH_SERVICE) as BluetoothManager
        bluetoothManager.adapter
    }
//    val bluetoothLeScanner: BluetoothLeScanner = bluetoothAdapter.bluetoothLeScanner
    private var bluetoothGatt: BluetoothGatt? = null
    private val REQ_CODE: Int = 10001
    private val OBDII_ADDR: String = "B8:27:EB:19:80:D8"

    private var stop_button: Button? = null;

    private var max_rpm: Float = 0.0f
    private var start_fuel: Float = 0.0f
    private var min_fuel: Float = 200.0f
    private var max_speed: Float = 0.0f
    private var max_coolant_temp: Float = 0.0f
    private var max_throttle_position: Float = 0.0f

    private val req_perms: Array<String> = if (android.os.Build.VERSION.SDK_INT >= android.os.Build.VERSION_CODES.S) {
        arrayOf(
            Manifest.permission.BLUETOOTH_SCAN,
            Manifest.permission.BLUETOOTH_CONNECT,
            Manifest.permission.ACCESS_FINE_LOCATION,
            Manifest.permission.BLUETOOTH
        )
    } else {
        arrayOf(Manifest.permission.ACCESS_FINE_LOCATION)
    }

    private val characteristics: Array<UUID> = arrayOf(
            UUID.fromString("000027af-0000-1000-8000-00805f9b34fb"), // rpm
            UUID.fromString("0000272f-0000-1000-8000-00805f9b34fb"), // coolant_temp
            UUID.fromString("00002730-0000-1000-8000-00805f9b34fb"), // intake_air_temp
            UUID.fromString("000027a7-0000-1000-8000-00805f9b34fb"), // speed
            UUID.fromString("00002731-0000-1000-8000-00805f9b34fb"), // ambient_temp
            UUID.fromString("000027ad-0000-1000-8000-00805f9b34fb"), // fuel_level
            UUID.fromString("000027c1-0000-1000-8000-00805f9b34fb"), // maf_flow_rate
            UUID.fromString("000027ae-0000-1000-8000-00805f9b34fb"), // throttle_pos
            UUID.fromString("00002b18-0000-1000-8000-00805f9b34fb"), // battery_voltage
            UUID.fromString("000027a4-0000-1000-8000-00805f9b34fb"), // odometer
            UUID.fromString("00002732-0000-1000-8000-00805f9b34fb"), // engine_oil_temp
            UUID.fromString("00002c08-0000-1000-8000-00805f9b34fb"), // gear_ratio
        )

    private val NOTIFICATION_DESCRIPTOR: UUID = UUID.fromString("00002902-0000-1000-8000-00805f9b34fb")

    private val operationQueue: Queue<Runnable> = ConcurrentLinkedQueue()
    private var isProcessingQueue = false

    // Add this method to enqueue operations
    private fun enqueueOperation(operation: Runnable) {
        operationQueue.add(operation)
        processQueue()
    }

    // Process the next operation in the queue
    private fun processQueue() {
        if (isProcessingQueue || operationQueue.isEmpty()) return

        isProcessingQueue = true
        val operation = operationQueue.poll()
        operation?.run()
    }

    // Call this when an operation is complete
    private fun operationCompleted() {
        isProcessingQueue = false
        processQueue()
    }


    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        window.decorView.apply {
            // Hide both the navigation bar and the status bar.
            // SYSTEM_UI_FLAG_FULLSCREEN is only available on Android 4.1 and higher, but as
            // a general rule, you should design your app to hide the status bar whenever you
            // hide the navigation bar.
            systemUiVisibility = View.SYSTEM_UI_FLAG_HIDE_NAVIGATION;
        }

        if (!(req_perms.all {ContextCompat.checkSelfPermission(this, it) == PackageManager.PERMISSION_GRANTED})) {
            ActivityCompat.requestPermissions(this, req_perms, REQ_CODE)
        } else {
            val device = bluetoothAdapter.getRemoteDevice(OBDII_ADDR)
            bluetoothGatt = device.connectGatt(this@MainActivity, false, gattCallback)
        }
        this.requestWindowFeature(Window.FEATURE_NO_TITLE);

        setContentView(R.layout.layout)

        stop_button = findViewById(R.id.stopButton);

        stop_button?.setOnClickListener {
            showMaxes();
        }
    }

//    private val leCallback = object: ScanCallback() {
//        override fun onScanResult(callbackType: Int, result: ScanResult?) {
//            super.onScanResult(callbackType, result)
//            if (result != null) {
//                Toast.makeText(this@MainActivity, result.device.name, Toast.LENGTH_LONG).show()
//            }
//        }
//    }

    override fun onRequestPermissionsResult(
        requestCode: Int,
        permissions: Array<out String>,
        grantResults: IntArray,
        deviceId: Int
    ) {
        super.onRequestPermissionsResult(requestCode, permissions, grantResults, deviceId)
        if (requestCode == REQ_CODE) {
            if (grantResults.all {it == PackageManager.PERMISSION_GRANTED}) {
                val device = bluetoothAdapter.getRemoteDevice(OBDII_ADDR)
                bluetoothGatt = device.connectGatt(this@MainActivity, false, gattCallback)
            }
        }
    }

    private val gattCallback = object: BluetoothGattCallback() {
        override fun onConnectionStateChange(gatt: BluetoothGatt?, status: Int, newState: Int) {
            val deviceAddress = gatt?.device?.address
            if (status == BluetoothGatt.GATT_SUCCESS) {
                if (newState == BluetoothProfile.STATE_CONNECTED) {
                    Log.w("BluetoothGattCallback", "Successfully connected to $deviceAddress")
                    gatt?.discoverServices()
                } else if (newState == BluetoothProfile.STATE_DISCONNECTED) {
                    Log.w("BluetoothGattCallback", "Successfully disconnected from $deviceAddress")
                    gatt?.close()
                }
            } else {
                Log.w("BluetoothGattCallback", "Error $status encountered for $deviceAddress! Disconnecting...")
                gatt?.close()
            }
        }

        override fun onServicesDiscovered(gatt: BluetoothGatt?, status: Int) {
            super.onServicesDiscovered(gatt, status)
            if (status == BluetoothGatt.GATT_SUCCESS) {
                val service = gatt?.getService(UUID.fromString("00001812-0000-1000-8000-00805f9b34fb"))
                characteristics.forEach {
                    val characteristic = service?.getCharacteristic(it)
                    if (characteristic != null) {
                        enqueueOperation( Runnable {
                            gatt?.setCharacteristicNotification(characteristic, true)
                            val descriptor = characteristic.getDescriptor(UUID.fromString("00002902-0000-1000-8000-00805f9b34fb"))
                            if (descriptor != null) {
                                Log.i("BLE", descriptor.toString())
                            } else {
                                Log.e("BLE", it.toString())
                            }
                            descriptor?.setValue(BluetoothGattDescriptor.ENABLE_NOTIFICATION_VALUE)
                            gatt?.writeDescriptor(descriptor)
                            Log.e("BLE", "Characterisitc ${characteristic.uuid.toString()})")
                        })
                    }
                }
            }
        }

        override fun onCharacteristicChanged(
            gatt: BluetoothGatt,
            characteristic: BluetoothGattCharacteristic,
            value: ByteArray
        ) {
            super.onCharacteristicChanged(gatt, characteristic, value)
            Log.i("BLE", "Characteristic ${characteristic.uuid} = ${value.contentToString()} = ${parseBytes(value)}")

            var text: TextView? = null
            var name: String = ""
            var units: String = ""
            var multiplier: Float = 1.0f
            val fval: Float = parseBytes(value)

            when (characteristic.uuid) {
                characteristics[0] -> {
                    name = "Engine Speed: "
                    units = " rpm"
                    text = findViewById(R.id.rpmView)
                    if (fval > max_rpm) {
                        max_rpm = fval
                    }
                }

                characteristics[1] -> {
                    name = "Coolant Temperature: "
                    units = " \u00b0C"
                    text = findViewById(R.id.coolantTempView)
                    if (fval > max_coolant_temp) {
                        max_coolant_temp = fval
                    }
                }

                characteristics[2] -> {
                    name = "Intake Air Temperature: "
                    units = " \u00b0C"
                    text = findViewById(R.id.intakeAirTempView)
                }

                characteristics[3] -> {
                    name = "Speed: "
                    units = " mph"
                    text = findViewById(R.id.speedView)
                    multiplier = 1.0f / 1.609f
                    if (fval > max_speed) {
                        max_speed = fval
                    }
                }

                characteristics[4] -> {
                    name = "Ambient Temperature: "
                    units = " \u00b0C"
                    text = findViewById(R.id.ambientTempView)
                }

                characteristics[5] -> {
                    name = "Fuel Level: "
                    units = "%"
                    text = findViewById(R.id.fuelView)
                    if (start_fuel == 0.0f) {
                        start_fuel = fval
                    }
                    if (fval < min_fuel) {
                        min_fuel = fval
                    }
                }

                characteristics[6] -> {
                    name = "Mass Air Flow Rate: "
                    units = " g/s"
                    text = findViewById(R.id.mafFlowRateView)
                }

                characteristics[7] -> {
                    name = "Throttle Position: "
                    units = "%"
                    text = findViewById(R.id.throttlePositionView)
                    if (fval > max_throttle_position) {
                        max_throttle_position = fval
                    }
                }

                characteristics[8] -> {
                    name = "Control Module Voltage: "
                    units = " V"
                    text = findViewById(R.id.voltageView)
                }

                characteristics[9] -> {
                    name = "Odometer: "
                    units = " miles"
                    text = findViewById(R.id.odometerView)
                    multiplier = 1.0f / 1.609f
                }

//                characteristics[10] -> {
//                    name = "Engine Oil Temperature: "
//                    units = " \u00b0C"
//                    text = findViewById(R.id.engineOilTempView)
//                }

//                characteristics[11] -> {
//                    name = "Gear Ratio: "
//                    units = ""
//                    text = findViewById(R.id.gearRatioView)
//                }
            }

            /*if (text != findViewById(R.id.gearRatioView)) {
                text?.post(Runnable {
                    text.text = String.format(Locale.getDefault(), "%s%.02f%s", name, fval * multiplier, units)
                })
            } else {*/
                text?.post(Runnable {
                    text.text = String.format(Locale.getDefault(), "%s%s", name, value.toString())
                })
            //}
        }

        override fun onDescriptorWrite(
            gatt: BluetoothGatt?,
            descriptor: BluetoothGattDescriptor?,
            status: Int
        ) {
            super.onDescriptorWrite(gatt, descriptor, status)
            operationCompleted()
        }
    }

    override fun onDestroy() {
        super.onDestroy()
        bluetoothGatt?.close()
        bluetoothGatt = null
    }

    private fun showMaxes() {
        var text: TextView? = null

        text = findViewById(R.id.rpmView)
        text?.text = String.format(Locale.getDefault(), "Max RPM: %.02fRPM", max_rpm)

        text = findViewById(R.id.throttlePositionView)
        text?.text = String.format(Locale.getDefault(), "Max Throttle Position: %.02f%s", max_throttle_position, "%")

        text = findViewById(R.id.speedView)
        text?.text = String.format(Locale.getDefault(), "Max Speed: %.02fMPH", max_speed / 1.609f)

        text = findViewById(R.id.coolantTempView)
        text?.text = String.format(Locale.getDefault(), "Max coolant temp: %.02f Degrees\u00b0C", max_coolant_temp)
    }
}

fun parseBytes(value: ByteArray): Float {
    if (value.size != 4) return 0.0F;
    return Float.fromBits(
        ((value[3].toInt() and 0xFF) shl 24)
        or ((value[2].toInt() and 0xFF) shl 16)
        or ((value[1].toInt() and 0xFF) shl 8)
        or (value[0].toInt() and 0xFF))
}
