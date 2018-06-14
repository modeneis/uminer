package drivers

import (
	"github.com/Sirupsen/logrus"

	"os"

	"github.com/modeneis/uminer/src/model"
	"github.com/modeneis/uminer/src/providers"
	"github.com/modeneis/uminer/src/providers/sia"
)

var (
	log *logrus.Entry = logrus.WithField("package", "drivers")
)

func init() {
	providers.UseProviders(
		sia.New(),
	)
}

// SetLogger set the logger
func SetLogger(loggers *logrus.Entry) {
	log = loggers.WithFields(log.Data)
}

// Run starts the miner workers
func Run(flags model.Flags, loggers *logrus.Entry) {

	SetLogger(loggers)

	coinType := flags.CoinType
	provider, err := providers.GetProvider(coinType)
	if err != nil {
		log.WithError(err).Error("Could not GetProvider")
		os.Exit(1)
	}

	//TODO: isolate and make this vendored
	err = provider.ConnectClient(&flags)
	if err != nil {
		log.WithError(err).Error("Could not connect to client")
		os.Exit(1)
	}

	err = provider.Start()
	if err != nil {
		log.WithError(err).Error("Could not start miner")
		os.Exit(1)
	}
}
