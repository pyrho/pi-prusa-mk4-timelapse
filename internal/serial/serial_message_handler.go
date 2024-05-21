package serial

import (
	"log"
	"strings"

	"github.com/pyrho/timelapse-serial/internal/camera"
	"github.com/pyrho/timelapse-serial/internal/ffmpeg"
)

const (
	COMMAND_CAPTURE = iota
	COMMAND_PRINT_START
	COMMAND_PRINT_STOP
	COMMAND_UNHANDLED
)

func CreateSerialMessageHandler(cam camera.CameraWrapperInterface, outputResolution *string) func(m string) {
	return func(message string) {

		switch command := parseCommand(message); command {

		case COMMAND_PRINT_START:
			log.Println("New print started")
			cam.Start()
			cam.CreateNewSnapshotsDir()
			log.Println("New photo directory created")

		case COMMAND_CAPTURE:
			log.Println("Capturing...")
			cam.Snap()

		case COMMAND_PRINT_STOP:
			log.Println("Print stopped")
			cam.Stop()
			log.Println("Print done, creating timelapse...")
			go ffmpeg.SpawnFFMPEG(cam.GetCurrentSnapshotsDir(), outputResolution)
		}

	}
}

func parseCommand(incomingMessage string) int {
	if strings.HasPrefix(incomingMessage, "// action:capture") {
		return COMMAND_CAPTURE
	} else if strings.HasPrefix(incomingMessage, "// status:print_start") {
		return COMMAND_PRINT_START
	} else if strings.HasPrefix(incomingMessage, "// status:print_stop") {
		return COMMAND_PRINT_STOP
	} else {
		return COMMAND_UNHANDLED
	}
}
