package web

import (
	"encoding/json"
	"io"
	"net/http"
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
