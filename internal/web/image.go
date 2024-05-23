package web

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"log"
	"math"
	"os"
	"regexp"

	// "strings"
	"github.com/daddye/vips"
)

func NewResize(imgPath string) []byte {
	// thumbPath := fmt.Sprintf("%s_thumb.jpg", imagePath)
	m1 := regexp.MustCompile(`snap([0-9]+.jpg)`)
	thumbPath := m1.ReplaceAllString(imgPath, "thumb${1}")
	thumb, err := os.Open(thumbPath)
	if err == nil {
		img, _, err := image.Decode(thumb)
		if err == nil {
			return imgToBytes(img)
		}
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
		fmt.Fprintln(os.Stderr, err)
		return []byte{}
	}

	//optional written to file
	err = os.WriteFile(thumbPath, buf, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return buf
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
