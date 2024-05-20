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
	instance           *gphoto2.Camera
	currentSnapshotDir string
	baseOutputDir      string
}
type CameraWrapperInterface interface {
	Start()
	Stop()
	Snap()
	SetSnapshotsDir(newOutputDir string)
}

func MakeCameraWrapper(baseOutputDir string) CameraWrapper {
	return CameraWrapper{baseOutputDir: baseOutputDir}
}

func (c *CameraWrapper) SetSnapshotsDir(newOutputDir string) {
	c.currentSnapshotDir = newOutputDir
}

func (c *CameraWrapper) Start() {
	c.instance = initCam()
	c.currentSnapshotDir = utils.CreateNewPhotoDirectory(c.baseOutputDir)
}

func (c *CameraWrapper) Stop() {
	if c.instance != nil {
		c.instance.Exit()
		c.instance.Free()
	}
}

func (c *CameraWrapper) Snap() {
	var currentSnapshotDir string
	if len(c.currentSnapshotDir) == 0 {
		currentSnapshotDir = c.baseOutputDir + "/orphans"
	} else {
		currentSnapshotDir = c.baseOutputDir + "/" + c.currentSnapshotDir
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
		panic(fmt.Sprintf("%s: %s", "Failed to connect to camera, make sure it's around!", err))
	}
	return c
}
