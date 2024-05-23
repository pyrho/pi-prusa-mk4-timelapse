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
	m1 := regexp.MustCompile(`capt([0-9]+.jpg)`)
	thumbPath := m1.ReplaceAllString(imgPath, "thumb${1}")
	// thumbPath := strings.Replace(imagePath, "/capt", "thumb", 1)
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

func ResizeMe(imagePath string) []byte {

	// thumbPath := fmt.Sprintf("%s_thumb.jpg", imagePath)
	m1 := regexp.MustCompile(`capt([0-9]+.jpg)`)
	thumbPath := m1.ReplaceAllString(imagePath, "thumb${1}")
	// thumbPath := strings.Replace(imagePath, "/capt", "thumb", 1)
	thumb, err := os.Open(thumbPath)
	if err == nil {
		img, _, err := image.Decode(thumb)
		if err == nil {
			return imgToBytes(img)
		}
	}

	f, err := os.Open(imagePath)
	if err != nil {
		log.Printf("Cannot open %s, %s", imagePath, err)
		return []byte{}
	}

	//encoding message is discarded, because OP wanted only jpg, else use encoding in resize function
	img, _, err := image.Decode(f)
	if err != nil {
		log.Println(err)
		return []byte{}
	}

	//this is the resized image
	resImg := resize(img, 750, 500)

	//this is the resized image []bytes
	imgBytes := imgToBytes(resImg)

	//optional written to file
	err = os.WriteFile(thumbPath, imgBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return imgBytes
}

func resize(img image.Image, length int, width int) image.Image {
	//truncate pixel size
	minX := img.Bounds().Min.X
	minY := img.Bounds().Min.Y
	maxX := img.Bounds().Max.X
	maxY := img.Bounds().Max.Y
	for (maxX-minX)%length != 0 {
		maxX--
	}
	for (maxY-minY)%width != 0 {
		maxY--
	}
	scaleX := (maxX - minX) / length
	scaleY := (maxY - minY) / width

	imgRect := image.Rect(0, 0, length, width)
	resImg := image.NewRGBA(imgRect)
	draw.Draw(resImg, resImg.Bounds(), &image.Uniform{C: color.Black}, image.Point{0, 0}, draw.Src)
	for y := 0; y < width; y += 1 {
		for x := 0; x < length; x += 1 {
			averageColor := getAverageColor(img, minX+x*scaleX, minX+(x+1)*scaleX, minY+y*scaleY, minY+(y+1)*scaleY)
			resImg.Set(x, y, averageColor)
		}
	}
	return resImg
}

func getAverageColor(img image.Image, minX int, maxX int, minY int, maxY int) color.Color {
	var averageRed float64
	var averageGreen float64
	var averageBlue float64
	var averageAlpha float64
	scale := 1.0 / float64((maxX-minX)*(maxY-minY))

	for i := minX; i < maxX; i++ {
		for k := minY; k < maxY; k++ {
			r, g, b, a := img.At(i, k).RGBA()
			averageRed += float64(r) * scale
			averageGreen += float64(g) * scale
			averageBlue += float64(b) * scale
			averageAlpha += float64(a) * scale
		}
	}

	averageRed = math.Sqrt(averageRed)
	averageGreen = math.Sqrt(averageGreen)
	averageBlue = math.Sqrt(averageBlue)
	averageAlpha = math.Sqrt(averageAlpha)

	averageColor := color.RGBA{
		R: uint8(averageRed),
		G: uint8(averageGreen),
		B: uint8(averageBlue),
		A: uint8(averageAlpha)}

	return averageColor
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
