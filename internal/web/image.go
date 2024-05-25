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

func exportAndWrite(image *vips.ImageRef, path string, ctx context.Context) error {
	select {
	case <-ctx.Done():
		log.Println("Goroutine closed by context cancel status")
		return nil
	default:
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

}

func resize(image *vips.ImageRef, ctx context.Context) error {
	select {
	case <-ctx.Done():
		log.Println("Goroutine closed by context cancel status")
		return nil
	default:
		if err := image.Resize(0.06, vips.KernelNearest); err != nil {
			return err
		}
		return nil
	}
}

func newImageFromFile(imgPath string, ctx context.Context) (*vips.ImageRef, error) {
	select {
	case <-ctx.Done():
		log.Println("Goroutine closed by context cancel status")
		return nil, nil
	default:
		image, err := vips.NewImageFromFile(imgPath)
		if err != nil {
			return nil, err
		}
		return image, nil
	}
}

func CreateAndSaveThumbnail(imgPath string, ctx context.Context) string {

	m1 := regexp.MustCompile(`snap([0-9]+.jpg)`)
	thumbPath := m1.ReplaceAllString(imgPath, "thumb${1}")
	if _, err := os.Stat(thumbPath); !errors.Is(err, fs.ErrNotExist) {
		// Thumbnail already exists
		return thumbPath
	}

	image, err := newImageFromFile(imgPath, ctx)
	if err != nil {
		log.Println("Cannot Open image from file", imgPath)
		return ""
	}
	defer image.Close()

	if err := resize(image, ctx); err != nil {
		log.Println("Cannot resize image", err)
		return ""
	}

	if err := exportAndWrite(image, thumbPath, ctx); err != nil {
		log.Println("Cannot export image", err)
		return ""

	}
	return thumbPath
}
