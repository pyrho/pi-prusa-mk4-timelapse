package web

import (
	"context"
	"errors"
	"io/fs"

	"log"
	"os"
	"regexp"

	"github.com/davidbyttow/govips/v2/vips"
)

func exportAndWrite(image *vips.ImageRef, path string) error {
	buf, _, err := image.ExportJpeg(&vips.JpegExportParams{
		Quality: 80,
	})

	if err != nil {
		return err
	}

	err = os.WriteFile(path, buf, 0644)
	if err != nil {
		return err
	}
	return nil
}

func resize(image *vips.ImageRef) error {
	if err := image.Resize(0.06, vips.KernelNearest); err != nil {
		return err
	}
	return nil
}

func newImageFromFile(imgPath string) (*vips.ImageRef, error) {
	image, err := vips.NewImageFromFile(imgPath)
	if err != nil {
		return nil, err
	}
	return image, nil
}

func CreateAndSaveThumbnail(imgPath string, ctx context.Context) string {
	select {
	case <-ctx.Done():
		log.Println("Thumbnail creation cancelled")
		return ""
	default:
		m1 := regexp.MustCompile(`snap([0-9]+.jpg)`)
		thumbPath := m1.ReplaceAllString(imgPath, "thumb${1}")
		if _, err := os.Stat(thumbPath); !errors.Is(err, fs.ErrNotExist) {
			// Thumbnail already exists
			return thumbPath
		}

		image, err := newImageFromFile(imgPath)
		if err != nil {
			log.Println("Cannot Open image from file", imgPath)
			return ""
		}
		defer func() {
			if image != nil {
				image.Close()
			}
		}()

		if err := resize(image); err != nil {
			log.Println("Cannot resize image", err)
			return ""
		}

		if err := exportAndWrite(image, thumbPath); err != nil {
			log.Println("Cannot export image", err)
			return ""

		}
		return thumbPath
	}

}
