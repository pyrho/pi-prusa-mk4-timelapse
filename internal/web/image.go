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
		return thumbPath
	}

	image, err := vips.NewImageFromFile(imgPath)
	if err != nil {
		panic(err)
	}
    defer image.Close()

	if err := image.Resize(0.06, vips.KernelLinear); err != nil {
		panic(err)
	}
	buf, _, err := image.ExportJpeg(&vips.JpegExportParams{
		Quality: 90,
	})
	err = os.WriteFile(thumbPath, buf, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return thumbPath
}
