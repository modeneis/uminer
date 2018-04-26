package util

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/robvanmieghem/go-opencl/cl"
)

// GetDevices get available devices to use on the platform target
func GetDevices(cpu int) (clDevices []*cl.Device) {

	var devicesTypesForMining = cl.DeviceTypeGPU
	if cpu > 0 {
		devicesTypesForMining = cl.DeviceTypeAll
	}

	platforms, err := cl.GetPlatforms()
	if err != nil {
		log.Panic("Error: ", err)
	}

	clDevices = make([]*cl.Device, 0, 4)
	for _, platform := range platforms {
		log.Println("INFO: Platform", platform.Name())
		platormDevices, err := cl.GetDevices(platform, devicesTypesForMining)
		if err != nil {
			log.Println("Error: ", err)
		}
		log.Println("INFO: ", len(platormDevices), "device(s) found:")
		for i, device := range platormDevices {
			log.Println("INFO: ", i, "-", device.Type(), "-", device.Name())
			clDevices = append(clDevices, device)
		}
	}

	if len(clDevices) == 0 {
		log.Println("ERROR: No suitable opencl devices found")
		os.Exit(1)
	}

	return clDevices
}

//DeviceExcludedForMining checks if the device is in the exclusion list
func DeviceExcludedForMining(deviceID int, excludedGPUs string) bool {
	excludedGPUList := strings.Split(excludedGPUs, ",")
	for _, excludedGPU := range excludedGPUList {
		if strconv.Itoa(deviceID) == excludedGPU {
			return true
		}
	}
	return false
}
