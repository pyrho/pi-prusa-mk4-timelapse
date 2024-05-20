package camera

import (
	"context"

	"github.com/rubiojr/go-usbmon"
)

func MonitorCameraUsbEvents(cameraSerialNumber *string, cameraWrapper *CameraWrapper) {

	actionFilter := &usbmon.ActionFilter{Action: usbmon.ActionAll}
	serialFilter := &usbmon.SerialFilter{Serial: *cameraSerialNumber}

	devs, err := usbmon.ListenFiltered(
		context.Background(),
		actionFilter,
		serialFilter,
	)

	if err != nil {
		panic(err)
	}

	for dev := range devs {
		switch dev.Action() {
		case "add":
			cameraWrapper.Start()
		case "remove":
			cameraWrapper.Stop()
		}
		// fmt.Printf("-- Device %s\n", dev.Action())
		// fmt.Println("Serial: " + dev.Serial())
		// fmt.Println("Path: " + dev.Path())
		// fmt.Println("Vendor: " + dev.Vendor())
	}
}
