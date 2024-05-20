package main

import (
	"github.com/pyrho/timelapse-serial/internal/camera"
	"github.com/pyrho/timelapse-serial/internal/config"
	"github.com/pyrho/timelapse-serial/internal/interrupt_trap"
	"github.com/pyrho/timelapse-serial/internal/serial"
)

type Command int

func main() {
	config := config.LoadConfig()
	camera := camera.MakeCameraWrapper(config.OutputDir)

	interrupttrap.TrapInterrupt(func() { camera.Stop() })

	serial.StartSerialRead(
		config.BaudRate,
		config.SerialPort,
		serial.CreateSerialMessageHandler(&camera),
	)
}
