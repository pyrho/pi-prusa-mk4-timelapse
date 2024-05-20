package ffmpeg

import (
	"fmt"
	"log"
	"os/exec"
)

func SpawnFFMPEG(capturedPhotosPath string) {
	// ffmpeg CMD: `ffmpeg -f image2 -framerate 24 -pattern_type glob -i "*.jpg" -crf 20 -c:v libx264 -pix_fmt yuv420p -s 1920x1280 output.mp4`
	log.Println("Starting FFMPEG timelapse creation at", capturedPhotosPath, "...")
	cmd := exec.Command(
		"ffmpeg",
		"-f", "image2", "-framerate", "24",
		"-pattern_type", "glob",
		"-i", fmt.Sprintf("%s/*.jpg", capturedPhotosPath),
		"-crf", "20",
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-s", "1920x1280",
		"-y",
		fmt.Sprintf("%s/output.mp4", capturedPhotosPath),
	)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Cannot create timelapse: %v", err)
	}
	log.Println("Timelapse created!")
}
