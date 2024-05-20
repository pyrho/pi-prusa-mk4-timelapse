package main

import (
	"flag"

	"github.com/pyrho/timelapse-serial/internal/camera"
	"github.com/pyrho/timelapse-serial/internal/interrupt_trap"
	"github.com/pyrho/timelapse-serial/internal/serial"
)

func main() {
	cameraSerial := flag.String("cameraSerial", "000007601060", "The serial number of the camera")
	portName := flag.String("portName", "/dev/ttyACM0", "The path of the printer port")
	baudRate := flag.Int("baudRate", 115200, "The baud rate of the serial port")
	outputDir := flag.String("outputDir", "/tmp/timelapse-serial-captures", "The output path where the pictures and timelapses will be stored")
	flag.Parse()

	c := camera.MakeCameraWrapper(*outputDir)
	camera.MonitorCameraUsbEvents(cameraSerial, &c)

	interrupttrap.TrapInterrupt(func() { c.Stop() })

	serial.StartSerialRead(
		*baudRate,
		*portName,
		serial.CreateSerialMessageHandler(&c),
	)
}
