package camera

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jonmol/gphoto2"
	"github.com/pyrho/timelapse-serial/internal/utils"
)

type CameraWrapper struct {
	instance                   *gphoto2.Camera
	currentSnapshotDirFullPath string
	baseOutputDir              string
}
type CameraWrapperInterface interface {
	Start()
	Stop()
	Snap()
	CreateNewSnapshotsDir()
}

func MakeCameraWrapper(baseOutputDir string) CameraWrapper {
	return CameraWrapper{baseOutputDir: baseOutputDir}
}

func (c *CameraWrapper) CreateNewSnapshotsDir() {
	c.currentSnapshotDirFullPath =
		utils.CreateNewPhotoDirectory(c.baseOutputDir)
	log.Println("Created new Snapshot directory: " + c.currentSnapshotDirFullPath)
}

func (c *CameraWrapper) Start() {
	if c.instance != nil {
		c.Stop()
	}
	c.instance = initCam()
	log.Println("Started CameraWrapper")
}

func (c *CameraWrapper) Stop() {
	if c.instance != nil {
		c.instance.Exit()
		c.instance.Free()
		c.instance = nil
		log.Println("Stopped cameraWrapper")
	} else {
		log.Println("Camera was already stopped")

	}
}

func (c *CameraWrapper) Snap() {
	var currentSnapshotDir string
	// This means that the program was spawned when a print
	// was already in progress.
	// We still want to save the pics, so just store them in the
	// orphans folder
	if len(c.currentSnapshotDirFullPath) == 0 {
		currentSnapshotDir = c.baseOutputDir + "/orphans"
	}

	snapFilename := fmt.Sprintf("%s/snap%d.jpg", currentSnapshotDir, time.Now().Unix())

	f, err := os.Create(snapFilename)
	if err != nil {
		log.Println("Failed to create temp file", snapFilename, "giving up!", err)
		return
	}

	if err := c.instance.CaptureDownload(f, false); err != nil {
		log.Println("Failed to capture!", err)
	}
}

func initCam() *gphoto2.Camera {
	// Calling `NewCamera` with `""` will connect to the first available camera
	c, err := gphoto2.NewCamera("")
	if err != nil {
        log.Println("No cameras detected")
        return nil
		// panic(fmt.Sprintf("%s: %s", "Failed to connect to camera, make sure it's around!", err))
	}
	return c
}
