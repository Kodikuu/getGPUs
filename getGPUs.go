package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func stringInArray(str string, arr []string) bool {
	for i, arrstr := range arr {
		_ = i
		if str == arrstr {
			return true
		}
	}
	return false
}

func filterEmpty(input []string) []string {
	var output []string

	for i, str := range input {
		_ = i
		if !stringInArray(str, []string{"", " "}) {
			output = append(output, str)
		}
	}
	return output
}

func getInfo() ([]string, [][]string) {
	rawdata, err := exec.Command("wmic", "path", "win32_videocontroller", "get", "name,pnpdeviceid,driverversion,driverdate").Output()
	if err != nil {
		log.Fatal(err)
	}

	rawstring := string(rawdata)
	rawlines := strings.Split(rawstring, "\r\r\n")

	rawheaders := rawlines[0]
	rawinfo := rawlines[1 : len(rawlines)-2]

	rawsplitheaders := strings.Split(rawheaders, "  ")
	var headers []string
	for i, s := range rawsplitheaders {
		_ = i
		if !stringInArray(s, []string{"", " "}) {
			headers = append(headers, strings.TrimSpace(s))
		}
	}

	var data [][]string
	for j, line := range rawinfo {
		_ = j
		rawsplitline := strings.Split(line, "  ")
		var dataline []string
		for i, str := range rawsplitline {
			_ = i
			dataline = append(dataline, strings.TrimSpace(str))
		}
		dataline = filterEmpty(dataline)
		data = append(data, dataline)
	}

	return headers, data
}

func dataToArray(headers []string, data [][]string) []map[string]string {
	var output []map[string]string

	for i, adapter := range data {
		_ = i
		adapterout := make(map[string]string)
		for j := range headers {
			adapterout[headers[j]] = adapter[j]
		}
		output = append(output, adapterout)
	}

	return output
}

func parseDevID(str string) (string, string) {
	vendor := str[8:12]
	device := str[17:21]

	return vendor, device
}

func parseDate(str string) []string {
	year := str[:4]
	month := str[4:6]
	day := str[6:8]
	output := []string{day, month, year}
	return output
}

func main() {

	headers, data := getInfo()
	datamaps := dataToArray(headers, data)

	for i, adapter := range datamaps {
		fmt.Printf("GPU %d:\n", i)
		for key, element := range adapter {
			if key == "PNPDeviceID" {
				vendor, device := parseDevID(element)
				fmt.Printf("Vendor ID: %s\n", vendor)
				fmt.Printf("Device ID: %s\n", device)
			} else if key == "DriverDate" {
				date := parseDate(element)
				fmt.Printf("DriverDate: %s/%s/%s\n", date[0], date[1], date[2])
			} else {
				fmt.Printf("%s: %s\n", key, element)
			}
		}
		fmt.Printf("\n")
	}

	fmt.Printf("\nPress Enter to terminate.")
	fmt.Scanln()
}