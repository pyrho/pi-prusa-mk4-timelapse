package web

import (
	"path/filepath"

	// "io/fs"
	"log"
	"os"

	// "os"
	// "path/filepath"
	// "github.com/pyrho/timelapse-serial/internal/utils"
	"regexp"
)

type SnapInfo struct {
	FolderName string
	FileName   string
	FilePath   string
}

func getSnapsForTimelapse(outputDir string, folderName string) []SnapInfo {
	validSnap := regexp.MustCompile(`^capt[0-9]+.jpg$`)
	var tl []SnapInfo
	files, err := os.ReadDir(filepath.Join(outputDir, folderName))
	if err != nil {
		log.Fatalf("Cannot read output dir: %s", err)
	}
	for _, file := range files {
		if !file.IsDir() && validSnap.MatchString(file.Name()) {
			tl = append(tl, SnapInfo{
				FilePath:   filepath.Join(outputDir, file.Name()),
				FolderName: folderName,
				FileName:   file.Name(),
			})
		}
	}
	return tl
}
func getTimelapseFolders3(outputDir string) []TLInfo {
	validDir := regexp.MustCompile(`^[0-9-]+$`)
	var tl []TLInfo
	// var tl2 map[string][]string
	files, err := os.ReadDir(outputDir)
	if err != nil {
		log.Fatalf("Cannot read output dir: %s", err)
	}
	for _, file := range files {
		if file.IsDir() && validDir.MatchString(file.Name()) {
			tl = append(tl, TLInfo{
				FolderPath: filepath.Join(outputDir, file.Name()),
				FolderName: file.Name(),
			})
		}
	}
	return tl
}
func getTimelapseFolders2(outputDir string) Timelapses {
	validDir := regexp.MustCompile(`^[0-9-]+$`)
	validSnap := regexp.MustCompile(`^capt[0-9]+.jpg$`)
	var tl Timelapses
	// var tl2 map[string][]string
	tl3 := make(map[string][]string)
	files, err := os.ReadDir(outputDir)
	if err != nil {
		log.Fatalf("Cannot read output dir: %s", err)
	}
	for _, file := range files {
		var ttt Timelapse

		if file.IsDir() && validDir.MatchString(file.Name()) {
			tl3["ahi"] = []string{}
			ttt = Timelapse{
				FolderName: file.Name(),
				FolderPath: filepath.Join(outputDir, file.Name()),
				Snaps:      []TimelapseSnap{},
			}

			filesInFolder, err := os.ReadDir(filepath.Join(outputDir, file.Name()))
			if err != nil {
				log.Fatalf("Cannot read output dir: %s", err)
			}
			for _, fileInFolder := range filesInFolder {
				if !fileInFolder.Type().IsDir() && validSnap.MatchString(fileInFolder.Name()) {
					ttt.Snaps = append(ttt.Snaps, TimelapseSnap{
						FileName: fileInFolder.Name(),
						FilePath: filepath.Join(outputDir, fileInFolder.Name()),
					})
				}
			}

			tl = append(tl, ttt)
		}
	}
	return tl
}

// func getTimelapseFolders(outputDir string) Timelapses {
// 	validDir := regexp.MustCompile(`^[0-9-]+$`)
// 	validSnap := regexp.MustCompile(`^capt[0-9]+.jpg$`)
// 	var tl Timelapses
// 	// var tl2 map[string][]string
// 	tl3 := make(map[string][]string)
// 	filepath.WalkDir(outputDir, func(path string, d fs.DirEntry, e error) error {
// 		if d.IsDir() && validDir.MatchString(d.Name()) {
// 			tl3["ahi"] = []string{}
// 			tl = append(tl, Timelapse{
// 				FolderName: d.Name(),
// 				Snaps:      []string{},
// 			})
// 		}
// 		log.Println(path)
// 		return e
// 	})
// 	return tl
// }
