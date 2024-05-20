package main

import (
	"flag"

	"github.com/pyrho/timelapse-serial/internal/camera"
	"github.com/pyrho/timelapse-serial/internal/interrupt_trap"
	"github.com/pyrho/timelapse-serial/internal/serial"
)

func main() {
	portName := flag.String("portName", "/dev/ttyACM0", "The path of the printer port")
	baudRate := flag.Int("baudRate", 115200, "The baud rate of the serial port")
	outputDir := flag.String("outputDir", "/tmp/timelapse-serial-captures", "The output path where the pictures and timelapses will be stored")
	flag.Parse()

    camera.Monit()

	camera := camera.MakeCameraWrapper(*outputDir)

	interrupttrap.TrapInterrupt(func() { camera.Stop() })

	serial.StartSerialRead(
		*baudRate,
		*portName,
		serial.CreateSerialMessageHandler(&camera),
	)
}
