package config

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	BaudRate   int
	SerialPort string
	OutputDir  string
}

func LoadConfig() Config {
	var err error

	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Cannot find user home directory: %v", err)
	}

	configPath := fmt.Sprintf("%s/.config/timelapse-serial.toml", homedir)
	configFileInBytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Config file not found, expected it at ~/.config/timelapse-serial.toml: %v", err)
	}
	configFileString := string(configFileInBytes)

	var conf Config
	_, err = toml.Decode(configFileString, &conf)

	return conf
}
