package ffmpeg

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"time"
    
    "github.com/pyrho/timelapse-serial/internal/config"
)

func SpawnFFMPEG(capturedPhotosPath string, ffmpegConfig config.FFMPEG) {
	// ch := make(chan int)
	ctx, cancel := context.WithTimeoutCause(context.Background(), 3*time.Minute, errors.New("Timed out while creating timelapse"))
	defer cancel()

	// ffmpeg CMD: `ffmpeg -f image2 -framerate 24 -pattern_type glob -i "*.jpg" -crf 20 -c:v libx264 -pix_fmt yuv420p -s 1920x1280 output.mp4`
	log.Println("Starting FFMPEG timelapse creation at", capturedPhotosPath, "...")
	cmd := exec.CommandContext(ctx,
		"ffmpeg",
		"-f", "image2",
        "-framerate", "24",
		"-pattern_type", "glob",
		"-i", fmt.Sprintf("%s/*.jpg", capturedPhotosPath),
		"-crf", "20",
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-s", ffmpegConfig.OutputVideoResolution, // "1920x1280",
		"-y",
		fmt.Sprintf("%s/output.mp4", capturedPhotosPath),
	)
	if err := cmd.Run(); err != nil {
		log.Println("Error: " + err.Error())
		// ch <- -1
	} else {
		log.Println("Timelapse created!")
		// ch <- 0
	}
}
