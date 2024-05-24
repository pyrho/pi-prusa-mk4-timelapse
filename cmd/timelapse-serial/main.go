package main

import (
	"flag"
	"log"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/pyrho/timelapse-serial/internal/camera"
	"github.com/pyrho/timelapse-serial/internal/config"
	"github.com/pyrho/timelapse-serial/internal/interrupt_trap"
	"github.com/pyrho/timelapse-serial/internal/serial"
	"github.com/pyrho/timelapse-serial/internal/web"
)

func main() {
	configPath := flag.String("configPath", "/usr/local/etc/timelapse-serial.toml", "The path of the config file")
	flag.Parse()
	config := config.LoadConfig(*configPath)

	c := camera.MakeCameraWrapper(config.Camera.OutputDir)

	if len(config.Camera.CameraSerialNumber) > 0 {
		go camera.MonitorCameraUsbEvents(&config.Camera.CameraSerialNumber, &c)
	} else {
		log.Println("Not monitoring camera plug events")
	}

	vips.Startup(nil)
	vips.LoggingSettings(nil, vips.LogLevelCritical)
	interrupttrap.TrapInterrupt(func() {
		c.Stop()
		vips.Shutdown()
	})

	onSerialMessageHandler := serial.CreateSerialMessageHandler(&c, &config.FFMPEG)
	// This needs to be last
	go serial.StartSerialLoop(&config, onSerialMessageHandler)

	go web.StartWebServer(&config)

	log.Println("Running...")

	// This will make the program run forever (unless interrupted/killed)
	select {}
}
