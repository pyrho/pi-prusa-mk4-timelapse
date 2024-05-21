package interrupttrap

import (
	"log"
	"os"
	"os/signal"
)

func TrapInterrupt(atInterrupt func()) {

	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt)
	go func() {
		log.Println("Monitoring interrupt...")
		<-interruptChannel
		log.Println("Interrupt trapped, freeing Camera")
		atInterrupt()
		os.Exit(0)
	}()
}
