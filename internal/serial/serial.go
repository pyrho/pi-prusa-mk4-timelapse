package serial

import (
	"log"

	"go.bug.st/serial"
)

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

func StartSerialRead(baudRate int, portName string, onRead OnRead) {
	// Open the serial port
	mode := &serial.Mode{
		BaudRate: baudRate,
	}
	port, err := serial.Open(portName, mode)

	if err != nil {
		log.Fatalf("Error opening serial port: %v", err)
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
			log.Printf("Error: %v", err)
			return
		}
	}
}
