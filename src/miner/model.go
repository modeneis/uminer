package miner

import (
	"github.com/Sirupsen/logrus"
	"github.com/robvanmieghem/go-opencl/cl"

	"github.com/modeneis/uminer/src/model"
)

var (
	log *logrus.Entry = logrus.WithField("package", "miner")
)

// Miner actually mines :-)
type Miner struct {
	ClDevices         map[int]*cl.Device
	HashRateReports   chan *model.HashRateReport
	miningWorkChannel chan *miningWork
	//Intensity defines the GlobalItemSize in a human friendly way, the GlobalItemSize = 2^Intensity
	Intensity      int
	GlobalItemSize int
	Client         model.Client
}

//miningWork is sent to the mining routines and defines what ranges should be searched for a matching nonce
type miningWork struct {
	Header []byte
	Offset int
	Job    interface{}
}

//singleDeviceMiner actually mines on 1 opencl device
type singleDeviceMiner struct {
	ClDevice          *cl.Device
	MinerID           int
	HashRateReports   chan *model.HashRateReport
	miningWorkChannel chan *miningWork
	//Intensity defines the GlobalItemSize in a human friendly way, the GlobalItemSize = 2^Intensity
	Intensity      int
	GlobalItemSize int
	Client         model.HeaderReporter
}
