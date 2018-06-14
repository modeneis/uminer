package util

import (
	"strconv"
	"strings"

	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/robvanmieghem/go-opencl/cl"
)

var (
	log *logrus.Entry = logrus.WithField("package", "util")
)

//CreateEmptyBuffer calls CreateEmptyBuffer on the supplied context and logs and panics if an error occurred
func CreateEmptyBuffer(ctx *cl.Context, flags cl.MemFlag, size int) (buffer *cl.MemObject) {
	buffer, err := ctx.CreateEmptyBuffer(flags, size)
	if err != nil {
		log.WithError(err).Error("Could not CreateEmptyBuffer")
	}
	return buffer
}

// GetDevices get available devices to use on the platform target
func GetDevices(cpu int) (clDevices []*cl.Device, err error) {

	var devicesTypesForMining = cl.DeviceTypeGPU
	if cpu > 0 {
		devicesTypesForMining = cl.DeviceTypeAll
	}

	platforms, err := cl.GetPlatforms()
	if err != nil {
		log.WithError(err).Error("Could not GetPlatforms")
		return nil, err
	}

	clDevices = make([]*cl.Device, 0, 4)
	for _, platform := range platforms {
		log.Println("INFO: Platform", platform.Name())
		platformDevices, err := cl.GetDevices(platform, devicesTypesForMining)
		if err != nil {
			log.WithError(err).WithField("platform", platform).WithField("devicesTypesForMining", devicesTypesForMining).Error("Could not get Devices to mine")
			return nil, err
		}
		log.WithField("platformDevices", platformDevices).Debugf("device(s) found: %d", len(platformDevices))

		for i, device := range platformDevices {
			log.WithField("i", i).WithField("device.Type", device.Type()).WithField("device.Name", device.Name()).Debug("Appending platform devices")
			clDevices = append(clDevices, device)
		}
	}

	if len(clDevices) == 0 {
		err = fmt.Errorf("ERROR: No suitable opencl devices found")
	}

	return clDevices, err
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
