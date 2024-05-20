package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

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
