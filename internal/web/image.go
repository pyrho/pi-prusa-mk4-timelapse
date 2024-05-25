package web

import (
	"errors"
	"io/fs"

	"log"
	"os"
	"regexp"

	"github.com/davidbyttow/govips/v2/vips"
)

func CreateAndSaveThumbnail(imgPath string) string {
	m1 := regexp.MustCompile(`snap([0-9]+.jpg)`)
	thumbPath := m1.ReplaceAllString(imgPath, "thumb${1}")
	if _, err := os.Stat(thumbPath); !errors.Is(err, fs.ErrNotExist) {
		// Thumbnail already exists
		return thumbPath
	}

	image, err := vips.NewImageFromFile(imgPath)
	if err != nil {
		log.Println("Cannot Open image from file", imgPath)
		return ""
	}
	defer image.Close()

	if err := image.Resize(0.06, vips.KernelLinear); err != nil {
		log.Println("Cannot resize image", err)
		return ""
	}
	buf, _, err := image.ExportJpeg(&vips.JpegExportParams{
		Quality: 80,
	})

	if err != nil {
		log.Println("Cannot export JPEG file", err)
		return ""
	}

	err = os.WriteFile(thumbPath, buf, 0644)
	if err != nil {
		log.Println("Cannot write JPEG file", err)
		return ""
	}
	return thumbPath
}
