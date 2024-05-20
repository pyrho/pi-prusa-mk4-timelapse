package interrupttrap

import (
	"fmt"
	"os"
	"os/signal"
)

func TrapInterrupt(atInterrupt func()) {

	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt)
	go func() {
		<-interruptChannel
		fmt.Println("Interrupt trapped, freeing Camera")
		atInterrupt()
		os.Exit(0)
	}()
}
