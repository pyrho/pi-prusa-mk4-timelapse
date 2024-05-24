package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Printer struct {
	PortName string
	BaudRate int
}

type Camera struct {
	CameraSerialNumber string
	OutputDir          string
	LiveFeedURL        string
}

type Web struct {
	ThumbnailCreationMaxGoroutines int
}

type FFMPEG struct {
	OutputVideoResolution string
	FramesPerSecond       string
	Codec                 string
	PixelFormat           string
	TimeoutInMinutes      int
}

func (f *FFMPEG) WithDefaults() FFMPEG {
	var conf FFMPEG = *f
	if len(conf.OutputVideoResolution) == 0 {
		conf.OutputVideoResolution = "3246x2158"
	}

	if len(conf.FramesPerSecond) == 0 {
		conf.FramesPerSecond = "24"
	}

	if len(conf.Codec) == 0 {
		conf.Codec = "libx264"
	}

	if len(conf.PixelFormat) == 0 {
		conf.PixelFormat = "yuv420p"
	}

	return conf
}

type Config struct {
	Printer Printer
	Camera  Camera
	FFMPEG  FFMPEG
	Web     Web
}

func LoadConfig(configPath string) Config {
	var conf Config
	_, err := toml.DecodeFile(configPath, &conf)
	if err != nil {
		log.Panicln("Cannot parse config file", err)
	}

	return conf
}
