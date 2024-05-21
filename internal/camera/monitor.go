package camera

import (
	"context"
	"log"

	"github.com/rubiojr/go-usbmon"
)

func MonitorCameraUsbEvents(cameraSerialNumber *string, cameraWrapper *CameraWrapper) {

	actionFilter := &usbmon.ActionFilter{Action: usbmon.ActionAdd}
	serialFilter := &usbmon.SerialFilter{Serial: *cameraSerialNumber}

	devs, err := usbmon.ListenFiltered(
		context.Background(),
		actionFilter,
		serialFilter,
	)

	if err != nil {
		panic(err)
	}

    log.Println("Monitoring camera with ID" + *cameraSerialNumber + "...")
	for dev := range devs {
		switch dev.Action() {
		case "add":
			log.Println("Camera connected")
			cameraWrapper.Start()
		}
	}
}

/*
Design Log

## 2024-05-20
When filtering on "serial", we only get the "add" event, because the "remove"
event does not have a vendor/serial attached to it (only path).
It's fine though, we really only care about when the camera connects, at which
point we need to refresh the gPhoto handle, and restarting the cameraWrapper
instance will clear the previous instance.
*/
