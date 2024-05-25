package web

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"regexp"

	"github.com/h2non/bimg"
)

func CreateThumbnail(imgPath string) (string, error) {
	thumbnailPath := thumbnailPathFromImagePath(imgPath)
	if thumbnailAlreadyExists(thumbnailPath) {
		return thumbnailPath, nil
	}
	if err := resizeImageAndSaveThumbnail(imgPath, thumbnailPath); err != nil {
		return "", err
	}
	return thumbnailPath, nil
}

func thumbnailPathFromImagePath(imgPath string) string {
	m1 := regexp.MustCompile(`snap([0-9]+.jpg)`)
	return m1.ReplaceAllString(imgPath, "thumb${1}")
}

func thumbnailAlreadyExists(thumbPath string) bool {
	_, err := os.Stat(thumbPath)
	return !errors.Is(err, fs.ErrNotExist)
}

func resizeImageAndSaveThumbnail(imgPath string, thumbPath string) error {
	buffer, err := bimg.Read(imgPath)
	if err != nil {
		log.Println(err)
	}

	shrunk, err := bimg.NewImage(buffer).Resize(600, 400)
	if err != nil {
		return err
	}

	if err := bimg.Write(thumbPath, shrunk); err == nil {
		return err
	}
	return nil
}
