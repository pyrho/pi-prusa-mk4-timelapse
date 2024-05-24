package web

import (
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"log"
	"net/http"

	"github.com/pyrho/timelapse-serial/internal/config"
	"github.com/pyrho/timelapse-serial/internal/web/assets"
)

type TimelapseSnap struct {
	FileName string
	FilePath string
}
type Timelapse struct {
	FolderName string
	FolderPath string
	Snaps      []TimelapseSnap
}
type Timelapses []Timelapse

type TLInfo struct {
	FolderName string
	FolderPath string
}
type Hi struct {
	B64     string
	ix      int
	ImgPath string
}

type TmplData struct {
	FolderName       string
	Timelapses       []TLInfo
	CurrentTimelapse []SnapInfo
	AllThumbs        []Hi
}

func StartWebServer(conf *config.Config) {
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		// http.ServeFile(w, r, "relative/path/to/favicon.ico")
		http.ServeFileFS(w, r, assets.FavIcon, "favicon.ico")

	})
	http.Handle("/serve/", http.StripPrefix("/serve/", http.FileServer(http.Dir(conf.Camera.OutputDir))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServerFS(assets.StyleCSS)))

	http.HandleFunc("/get-thumb/{folderName}/{fileName}", func(w http.ResponseWriter, r *http.Request) {
		thumb := CreateAndSaveThumbnail(filepath.Join(conf.Camera.OutputDir, r.PathValue("folderName"), r.PathValue("fileName")))
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(thumb)
	})

	http.HandleFunc("/get-file/{folderName}/{fileName}", func(w http.ResponseWriter, r *http.Request) {
		thumb := CreateAndSaveThumbnail(filepath.Join(conf.Camera.OutputDir, r.PathValue("folderName"), r.PathValue("fileName")))
		imgBase64Str := base64.StdEncoding.EncodeToString(thumb)
		io.WriteString(w, fmt.Sprintf("<img class='img-fluid' id='img-display' src='data:image/jpeg;base64,%s'/>", imgBase64Str))
	})

	http.HandleFunc("/clicked/{folderName}", func(w http.ResponseWriter, r *http.Request) {
		folderName := r.PathValue("folderName")
		// {{{

		mu := sync.Mutex{}
		var allThumbs []Hi
		snaps := getSnapsForTimelapse(conf.Camera.OutputDir, folderName)
		var wg sync.WaitGroup
		for ix, snap := range snaps {
			wg.Add(1)
			go func() {
				defer wg.Done()
				imgPath := filepath.Join(conf.Camera.OutputDir, snap.FolderName, snap.FileName)
				thumb := CreateAndSaveThumbnail(imgPath)
				imgBase64Str := base64.StdEncoding.EncodeToString(thumb)
				mu.Lock()
				allThumbs = append(allThumbs, Hi{
					B64:     imgBase64Str,
					ix:      ix,
					ImgPath: snap.FolderName + "/" + snap.FileName,
				})
				mu.Unlock()
			}()

		}
		wg.Wait()
		slices.SortFunc(allThumbs, func(a, b Hi) int {
			return a.ix - b.ix
		})
		///}}}
		timelapseVideoPath := fmt.Sprintf("%s/%s/output.mp4", conf.Camera.OutputDir, folderName)
		hasTimelapseVideo := true
		if _, err := os.Stat(timelapseVideoPath); errors.Is(err, os.ErrNotExist) {
			// path/to/whatever does not exist
			hasTimelapseVideo = false
		}
		template := template.Must(template.ParseFS(Templates, "templates/snaps.html"))
		template.ExecuteTemplate(w, "snaps", map[string]interface{}{
			"AllThumbs":    allThumbs,
			"FolderName":   folderName,
			"HasTimelapse": hasTimelapseVideo,
		})
	})

	http.HandleFunc("/modal/{folder}/{file}", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFS(Templates, "templates/modal.html"))
		template.ExecuteTemplate(w, "modal", map[string]interface{}{
			"ImgPath": r.PathValue("folder") + "/" + r.PathValue("file"),
		})
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		tl := getTimelapseFolders3(conf.Camera.OutputDir)
		// si := getSnapsForTimelapse(conf.Camera.OutputDir, tl[0].FolderName)

		// {{{

		mu := sync.Mutex{}
		var allThumbs []Hi
		snaps := getSnapsForTimelapse(conf.Camera.OutputDir, tl[0].FolderName)
		var wg sync.WaitGroup
		for ix, snap := range snaps {
			wg.Add(1)
			go func() {
				defer wg.Done()
				imgPath := filepath.Join(conf.Camera.OutputDir, snap.FolderName, snap.FileName)
				thumb := CreateAndSaveThumbnail(imgPath)
				imgBase64Str := base64.StdEncoding.EncodeToString(thumb)
				mu.Lock()
				allThumbs = append(allThumbs, Hi{
					B64:     imgBase64Str,
					ix:      ix,
					ImgPath: snap.FolderName + "/" + snap.FileName,
				})
				mu.Unlock()
			}()

		}
		wg.Wait()
		slices.SortFunc(allThumbs, func(a, b Hi) int {
			return a.ix - b.ix
		})
		///}}}
		timelapseVideoPath := fmt.Sprintf("%s/%s/output.mp4", conf.Camera.OutputDir, tl[0].FolderName)
		hasTimelapseVideo := true
		if _, err := os.Stat(timelapseVideoPath); errors.Is(err, os.ErrNotExist) {
			// path/to/whatever does not exist
			hasTimelapseVideo = false
		}
		templateData := map[string]interface{}{
			"Timelapses": tl,
			// "CurrentTimelapse": si,
			"AllThumbs":    allThumbs,
			"HasTimelapse": hasTimelapseVideo,
			"FolderName":   tl[0].FolderName,
			"LiveFeedURL":  conf.Camera.LiveFeedURL,
		}

		template := template.Must(template.ParseFS(Templates, "templates/layout.html", "templates/folders.html", "templates/snaps.html"))
		if err := template.Execute(w, templateData); err != nil {
			log.Fatal(err)
		}
	})
	log.Println("HTTP server running")
	log.Fatal(http.ListenAndServe(":3025", nil))

}
