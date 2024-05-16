package main

import (
	"fmt"
	"log"
	"go.bug.st/serial"
)

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
		case err := <-errChan:
			log.Printf("Error: %v", err)
			return
		}
	}
}

