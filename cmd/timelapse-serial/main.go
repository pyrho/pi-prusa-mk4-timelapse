package main

import (
	"flag"
	"log"

	"github.com/pyrho/timelapse-serial/internal/camera"
	"github.com/pyrho/timelapse-serial/internal/config"
	"github.com/pyrho/timelapse-serial/internal/interrupt_trap"
	"github.com/pyrho/timelapse-serial/internal/serial"
)

func main() {
	configPath := flag.String("configPath", "~/.config/timelapse-serial.toml", "The path of the config file")
	flag.Parse()
	config := config.LoadConfig(*configPath)

	c := camera.MakeCameraWrapper(config.Capture.OutputDir)

	if len(config.Camera.CameraSerialNumber) > 0 {
		go camera.MonitorCameraUsbEvents(&config.Camera.CameraSerialNumber, &c)
	} else {
		log.Println("Not monitoring camera plug events")
	}

	interrupttrap.TrapInterrupt(func() { c.Stop() })

	// This needs to be last
	serial.StartSerialRead(
		config.Printer.BaudRate,
		config.Printer.PortName,
		serial.CreateSerialMessageHandler(&c, &config.FFMPEG),
	)
}
