package web

import (
	"bytes"
	"errors"
	"io/fs"

	// "fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/daddye/vips"
)

func CreateAndSaveThumbnail(imgPath string) string {
	m1 := regexp.MustCompile(`snap([0-9]+.jpg)`)
	thumbPath := m1.ReplaceAllString(imgPath, "thumb${1}")
	if _, err := os.Stat(thumbPath); !errors.Is(err, fs.ErrNotExist) {
		return thumbPath
	}

	options := vips.Options{
		Width:        195,
		Height:       130,
		Crop:         false,
		Extend:       vips.EXTEND_WHITE,
		Interpolator: vips.BILINEAR,
		Gravity:      vips.CENTRE,
		Quality:      80,
	}
	f, _ := os.Open(imgPath)
	inBuf, _ := io.ReadAll(f)
	buf, err := vips.Resize(inBuf, options)
	if err != nil {
		// fmt.Fprintln(os.Stderr, err)
		return ""
	}

	//optional written to file
	err = os.WriteFile(thumbPath, buf, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return thumbPath
}

func imgToBytes(img image.Image) []byte {
	var opt jpeg.Options
	opt.Quality = 80

	buff := bytes.NewBuffer(nil)
	err := jpeg.Encode(buff, img, &opt)
	if err != nil {
		log.Println(err)
		return []byte{}
	}

	return buff.Bytes()
}
