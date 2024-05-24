package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"time"
)

func CreateDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func CreateNewPhotoDirectory(basePath string) string {
	newDirPath := fmt.Sprintf("%s/%s", basePath, time.Now().Format("2006-01-02-15-04-05"))
	// time.Now().Format("2006-01-02-15-04-05")
	if _, err := os.Stat(newDirPath); os.IsNotExist(err) {
		if err = os.MkdirAll(newDirPath, os.ModePerm); err != nil {
			log.Fatal("Cannot create directory", newDirPath)
		}
	}
	return newDirPath
}

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}
