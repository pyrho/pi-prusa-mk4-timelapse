package main

import (
    "os"
	"fmt"
	"log"
	"go.bug.st/serial"
    "github.com/jonmol/gphoto2"
)

func initCam() *gphoto2.Camera {
	camera, err := gphoto2.NewCamera("")
	if err != nil {
		panic(fmt.Sprintf("%s: %s", "Failed to connect to camera, make sure it's around!", err))
	}
    return camera
}

func snap(camera *gphoto2.Camera) {
    snapFile := "/tmp/testshot.jpeg"
	if f, err := os.Create(snapFile); err != nil {
		fmt.Println("Failed to create temp file", snapFile, "giving up!", err)
	} else {
		fmt.Println("Taking shot, then copy to", snapFile)
		if err := camera.CaptureDownload(f, false); err != nil {
			fmt.Println("Failed to capture!", err)
		}
	}
}

func readFromSerial(port serial.Port, dataChan chan<- string, errChan chan<- error) {
	buf := make([]byte, 100)
	for {
		n, err := port.Read(buf)
		if err != nil {
			errChan <- err
			return
		}
		if n > 0 {
			dataChan <- string(buf[:n])
		}
	}
}

func main() {
    camera := initCam()
	// Open the serial port
	mode := &serial.Mode{
		BaudRate: 9600,
	}
	port, err := serial.Open("/dev/ttyACM0", mode)
	if err != nil {
		log.Fatalf("Error opening serial port: %v", err)
	}
	defer port.Close()

	fmt.Println("Serial port opened successfully")

	// Create channels to handle data and errors
	dataChan := make(chan string)
	errChan := make(chan error)

	// Start a goroutine to read from the serial port
	go readFromSerial(port, dataChan, errChan)

	// Main loop to handle incoming data
	for {
		select {
		case data := <-dataChan:
			fmt.Printf("Received: %s\n", data)
            snap(camera)
		case err := <-errChan:
			log.Printf("Error: %v", err)
			return
		}
	}
}

