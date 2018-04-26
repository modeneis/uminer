package drivers

import (
	"github.com/robvanmieghem/go-opencl/cl"

	"github.com/Sirupsen/logrus"

	"github.com/modeneis/uminer/src/model"
	"github.com/modeneis/uminer/src/util"
)

var (
	log *logrus.Entry = logrus.WithField("package", "drivers")
)

// SetLogger set the logger
func SetLogger(loggers *logrus.Entry) {
	log = loggers.WithFields(log.Data)
}

// Run starts the miner workers
func Run(flags model.Flags, loggers *logrus.Entry) {

	SetLogger(loggers)

	clDevices := util.GetDevices(flags.CPU)

	//Filter the excluded devices
	miningDevices := make(map[int]*cl.Device)
	for i, device := range clDevices {
		if util.DeviceExcludedForMining(i, flags.ExcludeGPUS) {
			continue
		}
		miningDevices[i] = device
	}

}
