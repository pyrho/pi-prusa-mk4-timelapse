package serial

import (
	// "fmt"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pyrho/timelapse-serial/internal/config"
	"go.bug.st/serial"
)

const WAIT_TIME = 5 * time.Second

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

type OnRead func(s string)

func waitForSerialPort(portName string) error {
	for {
		_, err := os.Stat(portName)
		if err == nil {
			log.Println("Port '" + portName + "' exists")
			return nil
		}
		if !os.IsNotExist(err) {
			// continue
			return fmt.Errorf("failed to stat serial port: %v", err)
		}

		log.Println("Port '" + portName + "' is not ready yet, retrying")
		time.Sleep(WAIT_TIME)
	}
}

func StartSerialLoop(conf *config.Config, onRead OnRead) {

	for {
		if err := waitForSerialPort(conf.Printer.PortName); err != nil {
			log.Printf("Error: %v\n", err)
			return
		}

		log.Println("Serial port is now available")

		if err := openAndRead(conf, onRead); err != nil {
			log.Printf("Error: %v\n", err)
		}
		// Wait for a second before retrying the port
		time.Sleep(WAIT_TIME)
	}

}

func openAndRead(conf *config.Config, onRead OnRead) error {
	// Open the serial port
	mode := &serial.Mode{
		BaudRate: conf.Printer.BaudRate,
	}
	port, err := serial.Open(conf.Printer.PortName, mode)

	if err != nil {
		return fmt.Errorf("Error opening serial port: %v", err)
	}

	defer port.Close()

	log.Println("Serial port opened successfully")

	// Create channels to handle data and errors
	dataChan := make(chan string)
	errChan := make(chan error)

	log.Println("Ready...")
	// Start a goroutine to read from the serial port
	go readFromSerial(port, dataChan, errChan)

	// Main loop to handle incoming data
	for {
		select {
		case data := <-dataChan:
			// log.Printf("Received: %s", data)
			onRead(data)

		case err := <-errChan:
			return fmt.Errorf("Error: %v", err)
		}
	}
}
