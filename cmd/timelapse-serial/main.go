package main

import (
	"flag"
	"log"

	"github.com/pyrho/timelapse-serial/internal/camera"
	"github.com/pyrho/timelapse-serial/internal/interrupt_trap"
	"github.com/pyrho/timelapse-serial/internal/serial"
)

func main() {
	cameraSerialNumber := flag.String("cameraSerialNumber", "", "The serial number of the camera to monitor for connection events")
	portName := flag.String("portName", "/dev/ttyACM0", "The path of the printer port")
	baudRate := flag.Int("baudRate", 115200, "The baud rate of the serial port")
	outputDir := flag.String("outputDir", "/tmp/timelapse-serial-captures", "The output path where the pictures and timelapses will be stored")
    videoRes := flag.String("videoResolution", "3246x2158", "The output resolution of the timelapse video")
	flag.Parse()

	c := camera.MakeCameraWrapper(*outputDir)

	if len(*cameraSerialNumber) > 0 {
		go camera.MonitorCameraUsbEvents(cameraSerialNumber, &c)
	} else {
		log.Println("Not monitoring camera plug events")
	}

	interrupttrap.TrapInterrupt(func() { c.Stop() })

	// This needs to be last
	serial.StartSerialRead(
		*baudRate,
		*portName,
		serial.CreateSerialMessageHandler(&c, videoRes),
	)
}
