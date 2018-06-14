package sia

import (
	"math"

	"github.com/robvanmieghem/go-opencl/cl"

	"github.com/Sirupsen/logrus"

	"fmt"

	"github.com/modeneis/uminer/src/miner"
	"github.com/modeneis/uminer/src/model"
	"github.com/modeneis/uminer/src/util"
)

var (
	log *logrus.Entry = logrus.WithField("package", "sia")
)

// Provider is used only for testing.
type Provider struct {
	client model.Client
	flags  *model.Flags
}

// New creates a new fake SKY, and sets up important connection details.
func New() *Provider {
	//rest := &gui.Client{
	//	Addr: "https://explorer.skycoin.net" + ":" + "443" + "/api/",
	//}
	coinProvider := &Provider{}
	//coinFake.Start()
	return coinProvider
}

// Name is used only for testing.
func (p Provider) Name() string {
	return "sia"
}

// GetType is used to get TYPE eg:(SIA)
func (p Provider) GetType() string {
	return "SIA"
}

// ConnectClient is used to connect the client
func (p Provider) ConnectClient(fl *model.Flags) (err error) {

	p.flags = fl

	if p.flags.URL == "" || p.flags.Username == "" || p.flags.Password == "" {
		err = fmt.Errorf("Connect to Client must set URL, Username, Password")
		log.WithField("flags", fl).Error("ConnectClient with wrong flags")
	}

	return err
}

// Start is used to start the work
func (p Provider) Start() (err error) {

	clDevices, err := util.GetDevices(p.flags.CPU)
	if err != nil {
		log.WithError(err).Error("Could not Start, got error when running GetDevices")
		return err
	}

	//Filter the excluded devices
	miningDevices := make(map[int]*cl.Device)
	for i, device := range clDevices {
		if util.DeviceExcludedForMining(i, p.flags.ExcludeGPUS) {
			continue
		}
		miningDevices[i] = device
	}

	nrOfMiningDevices := len(miningDevices)
	var hashRateReportsChannel = make(chan *model.HashRateReport, nrOfMiningDevices*10)
	globalItemSize := int(math.Exp2(float64(p.flags.Intensity)))

	m := &miner.Miner{
		ClDevices:       miningDevices,
		HashRateReports: hashRateReportsChannel,
		Intensity:       p.flags.Intensity,
		GlobalItemSize:  globalItemSize,
		Client:          p.client,
	}
	m.Mine()

	return nil
}
