package web

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type printerInfo struct {
	State string
}
type jobInfo struct {
	Id            int
	Progress      float32
	TimeRemaining int `json:"time_remaining"`
	TimePrinting  int `json:"time_printing"`
}

type PrintInfo struct {
	Job     jobInfo
	Printer printerInfo
}

type PrintInfoCache struct {
	info PrintInfo
	mu   sync.RWMutex
}

func NewPrintInfoCache() *PrintInfoCache {
	return &PrintInfoCache{}
}

func (pi *PrintInfoCache) Get() PrintInfo {
	pi.mu.RLock()
	defer pi.mu.RUnlock()
	return pi.info
}

func (pi *PrintInfoCache) StartLoop(printerUrl, apiKey string) {
	ticker := time.NewTicker(10 * time.Second)
	// We never want to stop!
	// defer ticker.Stop()

	go func() {
		info, _ := getPrinterInformation(printerUrl, apiKey)
		pi.mu.Lock()
		pi.info = info
		pi.mu.Unlock()

		for range ticker.C {
			info, err := getPrinterInformation(printerUrl, apiKey)
			if err != nil {
				log.Printf("Cannot get printer info: %v\n", err)
			}
			pi.mu.Lock()
			pi.info = info
			pi.mu.Unlock()
		}
	}()
}

func getPrinterInformation(printeUrl string, apiKey string) (PrintInfo, error) {
	// Create HTTP client
	client := http.Client{}

	req, err := http.NewRequest("GET", printeUrl+"/api/v1/status", nil)
	if err != nil {
		return PrintInfo{}, err
	}
	req.Header.Set("X-Api-Key", apiKey)

	// Send request and get response
	resp, err := client.Do(req)
	if err != nil {
		return PrintInfo{}, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return PrintInfo{}, err
	}

	var printInf PrintInfo
	if err = json.Unmarshal(b, &printInf); err != nil {
		return PrintInfo{}, err
	}
	return printInf, nil

}

/*
Example output:
{
  "job": {
    "id": 650,
    "progress": 68.00,
    "time_remaining": 1680,
    "time_printing": 4157
  },
  "storage": {
    "path": "/usb/",
    "name": "usb",
    "read_only": false
  },
  "printer": {
    "state": "PRINTING",
    "temp_bed": 82.0,
    "target_bed": 82.0,
    "temp_nozzle": 192.0,
    "target_nozzle": 192.0,
    "axis_z": 72.5,
    "flow": 100,
    "speed": 100,
    "fan_hotend": 7999,
    "fan_print": 6064
  }
}
*/
