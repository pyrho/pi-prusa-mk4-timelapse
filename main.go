package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/jonmol/gphoto2"
	"go.bug.st/serial"
)

type config struct {
	BaudRate   int
	SerialPort string
	OutputDir  string
}

func spawnFFMPEG(capturedPhotosPath string) {
	// ffmpeg CMD: `ffmpeg -f image2 -framerate 24 -pattern_type glob -i "*.jpg" -crf 20 -c:v libx264 -pix_fmt yuv420p -s 1920x1280 output.mp4`
	log.Println("Starting FFMPEG timelapse creation at", capturedPhotosPath, "...")
	cmd := exec.Command(
		"ffmpeg",
		"-f", "image2", "-framerate", "24",
		"-pattern_type", "glob",
		"-i", fmt.Sprintf("%s/*.jpg", capturedPhotosPath),
		"-crf", "20",
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-s", "1920x1280",
		"-y",
		fmt.Sprintf("%s/output.mp4", capturedPhotosPath),
	)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Cannot create timelapse: %v", err)
	}
	log.Println("Timelapse created!")
}

func loadConfig() config {
	var err error

	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Cannot find user home directory: %v", err)
	}

	configPath := fmt.Sprintf("%s/.config/timelapse-serial.toml", homedir)
	configFileInBytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Config file not found, expected it at ~/.config/timelapse-serial.toml: %v", err)
	}
	configFileString := string(configFileInBytes)

	var conf config
	_, err = toml.Decode(configFileString, &conf)

	return conf
}

func initCam() *gphoto2.Camera {
	c, err := gphoto2.NewCamera("")
	if err != nil {
		panic(fmt.Sprintf("%s: %s", "Failed to connect to camera, make sure it's around!", err))
	}
	return c
}

func createNewPhotoDirectory(basePath string) string {
	newDirPath := fmt.Sprintf("%s/%s", basePath, time.Now().Format("2006-01-02-15-04-05"))
	// time.Now().Format("2006-01-02-15-04-05")
	if _, err := os.Stat(newDirPath); os.IsNotExist(err) {
		if err = os.MkdirAll(newDirPath, os.ModePerm); err != nil {
			log.Fatal("Cannot create directory %s: %v", newDirPath, err)
		}
	}
	return newDirPath
}

func snap(camera *gphoto2.Camera, path string) {
	if len(path) == 0 {
		log.Fatal("There is no folder created for this print, should not happen")
	}

	snapFile := fmt.Sprintf("%s/capt%d.jpg", path, time.Now().Unix())
	if f, err := os.Create(snapFile); err != nil {
		log.Println("Failed to create temp file", snapFile, "giving up!", err)
	} else {
		if err := camera.CaptureDownload(f, false); err != nil {
			log.Println("Failed to capture!", err)
		}
	}
}

func readFromSerial(port serial.Port, dataChan chan<- string, errChan chan<- error) {
	buf := make([]byte, 500)
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

var camera *gphoto2.Camera

func main() {
	var capturePath string
	camera = initCam()
	config := loadConfig()
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt)

	go func() {
		<-interruptChannel
		fmt.Println("Interrupt trapped, freeing Camera")
		camera.Exit()
		camera.Free()
		os.Exit(0)
	}()

	// Open the serial port
	mode := &serial.Mode{
		BaudRate: config.BaudRate,
	}
	port, err := serial.Open(config.SerialPort, mode)

	if err != nil {
		log.Fatalf("Error opening serial port: %v", err)
	}

	defer port.Close()

	log.Println("Serial port opened successfully")

	// Create channels to handle data and errors
	dataChan := make(chan string)
	errChan := make(chan error)

	// Start a goroutine to read from the serial port
	go readFromSerial(port, dataChan, errChan)

	// Main loop to handle incoming data
	for {
		select {
		case data := <-dataChan:
            log.Printf("Received: %s\n", data)

			if strings.Contains(data, "status:print_start") {
                log.Println("New print started, creating folder.")
				capturePath = createNewPhotoDirectory(config.OutputDir)
			}

			if strings.Contains(data, "action:capture") {
				log.Println("Capturing.")
				go snap(camera, capturePath)
			}

			if strings.Contains(data, "status:print_stop") {
				log.Println("Print done, creating timelapse")
				go spawnFFMPEG(capturePath)
			}
		case err := <-errChan:
			log.Printf("Error: %v", err)
			return
		}
	}
}
